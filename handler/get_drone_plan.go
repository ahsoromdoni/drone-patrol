package handler

import (
	"math"
	"net/http"

	"github.com/ahsoromdoni/drone-patrol/constant"
	"github.com/ahsoromdoni/drone-patrol/generated"
	"github.com/ahsoromdoni/drone-patrol/repository"
	"github.com/ahsoromdoni/drone-patrol/utils"
	"github.com/labstack/echo/v4"
)

type Drone struct {
	ToleranceDistance int
	CurrentHeight     *int
	LastHeight        *int
}

type Coordinate struct {
	X, Y int
}

// Get drone patrol plan for the estate
// (GET /estate/{id}/drone-plan)
func (s *Server) GetDronePlan(ctx echo.Context, id string, params generated.GetDronePlanParams) error {
	var getEstateByIdInput = repository.GetEstateByIdInput{Id: id}
	estate, err := s.Repository.GetEstateById(ctx.Request().Context(), getEstateByIdInput)
	if err != nil {
		return utils.CheckForNotFoundError(ctx, err, constant.EstateNotFound)
	}

	var getTreesByEstateId = repository.GetTreesByEstateIdInput{EstateId: id}
	trees, err := s.Repository.GetTreesByEstateId(ctx.Request().Context(), getTreesByEstateId)
	if err != nil {
		return utils.CheckForNotFoundError(ctx, err, constant.EstateNotFound)
	}

	var distance = getDroneTravelDistance(estate, trees)

	var resp = generated.DronePlanResponse{Distance: distance}
	return ctx.JSON(http.StatusOK, resp)
}

func generateTreeCoordinateMap(trees []repository.GetTreesByEstateIdOutput) map[Coordinate]int {
	treeMap := make(map[Coordinate]int)
	for _, tree := range trees {
		treeMap[Coordinate{tree.XAxis, tree.YAxis}] = tree.Height
	}
	return treeMap
}

func calculateVerticalMovement(tree int, drone Drone) (verticalMove int) {
	if *drone.CurrentHeight <= tree {
		*drone.CurrentHeight = tree + drone.ToleranceDistance
		verticalMove = int(math.Abs(float64(*drone.LastHeight - *drone.CurrentHeight)))
	} else {
		*drone.CurrentHeight = *drone.CurrentHeight - tree + drone.ToleranceDistance
		verticalMove = int(math.Abs(float64(*drone.LastHeight - *drone.CurrentHeight)))
	}

	*drone.LastHeight = *drone.CurrentHeight
	return verticalMove
}

func getDroneTravelDistance(estate repository.GetEstateByIdOutput, trees []repository.GetTreesByEstateIdOutput) int {
	var distance = 0
	var horizontalMove = 10
	var verticalMove = 0
	var lastDroneHeight = 0
	var currentDroneHeight = 0
	var distanceAdjustment = 10
	var isThereATree = false

	var drone = Drone{
		ToleranceDistance: 1,
		LastHeight:        &lastDroneHeight,
		CurrentHeight:     &currentDroneHeight,
	}

	var treeMap = generateTreeCoordinateMap(trees)

	for i := 1; i <= estate.Width; i++ {
		if i%2 != 0 {
			// Drone Fly to the tight
			for j := 1; j <= estate.Length; j++ {
				if tree, ok := treeMap[Coordinate{j, i}]; ok {
					distance += calculateVerticalMovement(tree, drone)
					isThereATree = true
				}

				distance += horizontalMove
			}
		} else {
			// Drone Fly to the left
			for j := estate.Length; j >= 1; j-- {
				if tree, ok := treeMap[Coordinate{j, i}]; ok {
					distance += calculateVerticalMovement(tree, drone)
					isThereATree = true
				}

				distance += horizontalMove
			}
		}
	}

	// If there is a tree in the drone's path, adjust vertical movement based on the current height.
	if isThereATree {
		verticalMove = *drone.CurrentHeight
	} else {
		// If there are no trees, give extra vertical distance based on tolerance.
		verticalMove = drone.ToleranceDistance * 2
	}

	// Add the necessary vertical movement to the total distance.
	distance += verticalMove

	// Apply the final adjustment to the distance, reducing any unnecessary distance.
	distance -= distanceAdjustment

	return distance
}
