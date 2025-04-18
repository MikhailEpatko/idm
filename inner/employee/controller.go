package employee

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"idm/inner/common"
	"idm/inner/web"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
}

// интерфейс сервиса employee.Service
type Svc interface {
	FindById(id int64) (Response, error)
	CreateEmployee(request CreateRequest) (int64, error)
}

func NewController(server *web.Server, employeeService Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

// функция для регистрации маршрутов
func (c *Controller) RegisterRoutes() {

	// полный маршрут получится "/api/v1/employees"
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
}

// функция-хендлер, которая будет вызываться при POST запросе по маршруту "/api/v1/employees"
func (c *Controller) CreateEmployee(ctx fiber.Ctx) error {

	// анмаршалим JSON body запроса в структуру CreateRequest
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// вызываем метод CreateEmployee сервиса employee.Service
	var newEmployeeId, err = c.employeeService.CreateEmployee(request)
	if err != nil {
		switch {

		// если сервис возвращает ошибку RequestValidationError или AlreadyExistsError,
		// то мы возвращаем ответ с кодом 400 (BadRequest)
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())

		// если сервис возвращает другую ошибку, то мы возвращаем ответ с кодом 500 (InternalServerError)
		default:
			return common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
	}

	// функция OkResponse() формирует и направляет ответ в случае успеха
	if err = common.OkResponse(ctx, newEmployeeId); err != nil {

		// функция ErrorResponse() формирует и направляет ответ в случае ошибки
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}
	return nil
}
