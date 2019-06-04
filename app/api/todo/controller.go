package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/app/api/common"
)

// Controller handles http request.
type Controller struct {
	repo                  Repository
	modelValidatorFactory func() TodoModelValidator
}

// NewController return new todo controller instance.
func NewController(repo Repository) *Controller {
	return &Controller{
		repo: repo,
		modelValidatorFactory: func() TodoModelValidator {
			return newTodoModelValidator()
		},
	}
}

// RegisterRoutes register handler routes.
func (controller Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", "/", controller.createTodo)
	router.Handle("GET", "/", controller.getAllTodos)
	router.Handle("GET", "/:id", controller.getTodoByTodoID)
	router.Handle("PUT", "/:id", controller.updateTodoByTodoID)
	router.Handle("DELETE", "/:id", controller.removeTodoByTodoID)
}

func (controller *Controller) createTodo(ctx *gin.Context) {
	todoModelValidator := controller.modelValidatorFactory()
	if err := todoModelValidator.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewError("error", err))
		return
	}

	todo := todoModelValidator.Todo()
	createdTodo, err := controller.repo.createTodo(todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(201, createdTodo)
}

func (controller *Controller) getAllTodos(ctx *gin.Context) {
	todos, err := controller.repo.getTodos()
	if err != nil {
		ctx.String(500, err.Error())
	}

	ctx.JSON(200, todos)
}

func (controller *Controller) getTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	todo, err := controller.repo.getTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.AbortWithStatus(404)
	} else if err != nil {
		ctx.AbortWithStatusJSON(500, map[string]string{
			"error": err.Error(),
		})
	}

	ctx.JSON(200, todo)
}

func (controller *Controller) updateTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")
	todo := Todo{}

	if err := ctx.BindJSON(&todo); err != nil {
		ctx.AbortWithError(400, err)
	}

	todo, err := controller.repo.updateTodoByTodoID(todoID, todo)
	if err == common.ErrEntityNotFound {
		ctx.AbortWithError(404, err)
	} else if err != nil {
		ctx.AbortWithError(500, err)
	}

	ctx.JSON(200, todo)
}

func (controller *Controller) removeTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	removedTodoID, err := controller.repo.removeTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.AbortWithError(404, err)
	} else if err != nil {
		ctx.AbortWithError(500, err)
	}

	ctx.JSON(204, removedTodoID)
}