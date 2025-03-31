package service

import (
	"music-lib/internal/infrastructure/email"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"
	"music-lib/pkg/event"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository repository.IUserRepository
	Event          *event.EventBus
}

func NewAuthService(
	userRepository repository.IUserRepository,
	event *event.EventBus,
) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		Event:          event,
	}
}

func (service *AuthService) Login(email, password string) (*model.User, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(repository.EmailKey, email)
	if existedUser == nil {
		return nil, er.ErrWrongUserCredentials
	}

	// Проверка, верифицирован ли пользователь
	if !existedUser.IsVerified {
		return nil, er.ErrUserNotVerified
	}

	// Сравниваем пароли
	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return nil, er.ErrWrongUserCredentials
	}

	return existedUser, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	// Находим пользователя и проверяем его наличие
	existedUser, _ := service.UserRepository.FindByKey(repository.EmailKey, email)

	if existedUser != nil && existedUser.IsVerified { // если пользователь существует и верифицирован
		return "", er.ErrUserExists
	} else if existedUser != nil { // если пользователь существует и НЕ верифицирован
		// Регенерация кода и id сессии
		existedUser.Generate()

		// Отправить сообщение с кодом
		service.sendEmail(email, existedUser.Code)

		// Перезаписать юзера в бд
		_, err := service.UserRepository.Update(existedUser)
		if err != nil {
			return "", err
		}

		return existedUser.SessionId, nil
	}

	// Генерим хеш пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Создаем модель юзера
	user := &model.User{
		Email:      email,
		Password:   string(hashedPassword),
		Name:       name,
		Role:       model.RoleUser,
		IsVerified: false,
	}
	user.Generate()

	// Создаем запись user
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}

	// Отправить сообщение с кодом
	service.sendEmail(email, user.Code)

	return user.SessionId, nil
}

func (service *AuthService) Verify(sessionId, code string) (*model.User, error) {
	// Находим пользователя
	existedUser, _ := service.UserRepository.FindByKey(repository.SessionIdKey, sessionId)
	if existedUser == nil {
		return nil, er.ErrUserExists
	}

	// Проверка на подлинность кода
	if existedUser.Code != code {
		return nil, er.ErrNotAuthorized
	}

	// Пользователь становится верифицированным
	existedUser.IsVerified = true
	user, err := service.UserRepository.Update(existedUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *AuthService) sendEmail(mail, code string) {
	// Отправка кода на почту
	go service.Event.Publish(event.Event{
		Type: event.EventSendEmail,
		Data: email.Addressee{
			To:      mail,
			Subject: "Подтвердите почту",
			Text:    "Ваш персональный код подтверждения личности: " + code + ". Не сообщайте никому данный код.",
		},
	})
}
