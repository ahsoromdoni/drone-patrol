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

	var maxDistance = 0
	if params.MaxDistance != nil {
		maxDistance = *params.MaxDistance
	}

	var distance, coordinate = getDroneTravelDistance(estate, trees, maxDistance)

	var resp interface{}
	if maxDistance > 0 {
		resp = generated.DronePlanResponseWithRest{
			Distance: distance,
			Rest: &generated.Rest{
				X: coordinate.X,
				Y: coordinate.Y,
			},
		}
	} else {
		resp = generated.DronePlanResponse{
			Distance: distance,
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func getDroneTravelDistance(estate repository.GetEstateByIdOutput, trees []repository.GetTreesByEstateIdOutput, maxDistance int) (int, Coordinate) {
	var distance = 0                                     // Variable to keep track of the total distance traveled by the drone
	var lastDroneHeight = 0                              // The height of the drone from the previous step
	var currentDroneHeight = 0                           // The current height of the drone
	var horizontalMove = constant.HorizontalMove         // The fixed distance the drone moves horizontally per step
	var distanceAdjustment = constant.DistanceAdjustment // Distance adjustment for vertical movement
	var isThereATree = false                             // Flag to indicate if the drone encountered a tree in the current position

	var drone = Drone{
		ToleranceDistance: constant.ToleranceDistance,
		LastHeight:        &lastDroneHeight,
		CurrentHeight:     &currentDroneHeight,
	}

	var treeMap = generateTreeCoordinateMap(trees)

	for y := 1; y <= estate.Width; y++ {
		if y%2 != 0 {
			// Drone fly to the right
			for x := 1; x <= estate.Length; x++ {
				// Check if there is a tree at the current coordinate
				if tree, ok := treeMap[Coordinate{x, y}]; ok {
					// If there is a tree, calculate the vertical movement
					distance += calculateVerticalMovement(tree, drone)
					isThereATree = true // Set the flag indicating a tree was encountered
				}

				distance += horizontalMove

				// Adjust the distance based on whether a tree was encountered
				adjustDistance := adjustDistance(isThereATree, distance, distanceAdjustment, drone)

				// Check if the drone has exceeded the maximum allowed distance
				distanceCheck, coordinate := checkMaxDistance(maxDistance, horizontalMove, adjustDistance, x, y, estate, true)
				if distanceCheck > 0 {
					return distanceCheck, coordinate
				}
			}
		} else {
			// Drone fly to the left
			for x := estate.Length; x >= 1; x-- {
				if tree, ok := treeMap[Coordinate{x, y}]; ok {
					distance += calculateVerticalMovement(tree, drone)
					isThereATree = true
				}

				distance += horizontalMove

				adjustDistance := adjustDistance(isThereATree, distance, distanceAdjustment, drone)
				distanceCheck, coordinate := checkMaxDistance(maxDistance, horizontalMove, adjustDistance, x, y, estate, false)
				if distanceCheck > 0 {
					return distanceCheck, coordinate
				}

			}
		}
	}

	// Final distance adjustment after all movement calculations
	distance = adjustDistance(isThereATree, distance, distanceAdjustment, drone)

	return distance, Coordinate{}
}

func generateTreeCoordinateMap(trees []repository.GetTreesByEstateIdOutput) map[Coordinate]int {
	treeMap := make(map[Coordinate]int)
	for _, tree := range trees {
		treeMap[Coordinate{tree.XAxis, tree.YAxis}] = tree.Height
	}
	return treeMap
}

func calculateVerticalMovement(tree int, drone Drone) (verticalMove int) {
	if *drone.CurrentHeight != (tree + 1) {
		if *drone.CurrentHeight <= tree {
			*drone.CurrentHeight = tree + drone.ToleranceDistance
			verticalMove = int(math.Abs(float64(*drone.LastHeight - *drone.CurrentHeight)))
		} else {
			*drone.CurrentHeight = *drone.CurrentHeight - tree + drone.ToleranceDistance
			verticalMove = int(math.Abs(float64(*drone.LastHeight - *drone.CurrentHeight)))
		}
	}

	*drone.LastHeight = *drone.CurrentHeight
	return verticalMove
}

func adjustDistance(isThereATree bool, distance, distanceAdjustment int, drone Drone) int {
	var verticalMove int

	if isThereATree {
		// If there is a tree in the drone's path, adjust vertical movement based on the current height.
		verticalMove = *drone.CurrentHeight
	} else {
		// If there are no trees, give extra vertical distance based on tolerance.
		verticalMove = drone.ToleranceDistance * 2
	}

	distance += verticalMove       // Add the necessary vertical movement to the total distance.
	distance -= distanceAdjustment // Apply the final adjustment to the distance, reducing any unnecessary distance.

	return distance
}

func checkMaxDistance(maxDistance, horizontalMove, distance, x, y int, estate repository.GetEstateByIdOutput, right bool) (int, Coordinate) {
	// If the maximum distance is zero or less, return zero and an empty coordinate.
	if maxDistance <= 0 {
		return 0, Coordinate{}
	}

	// If the maxDistance is less than or equal to the horizontalMove or
	// the total distance traveled matches maxDistance, return the current position.
	if maxDistance <= horizontalMove || distance == maxDistance {
		return maxDistance, Coordinate{X: x, Y: y}
	}

	if right { // If the drone is moving to the right
		// If maxDistance is greater than the current distance and the drone reaches the estate's end boundary,
		// return the end boundary coordinates.
		if maxDistance > distance && x == estate.Length && y == estate.Width {
			return distance, Coordinate{X: estate.Length, Y: estate.Width}
		}
		// If the distance exceeds maxDistance and the drone is at the leftmost edge,
		// move one row upward and return the new position.
		if distance > maxDistance && x == 1 {
			return maxDistance, Coordinate{X: x, Y: y - 1}
		}
		// If the distance exceeds maxDistance, move one step to the left and return the new position.
		if distance > maxDistance {
			return maxDistance, Coordinate{X: x - 1, Y: y}
		}
	} else { // If the drone is moving to the left
		// If maxDistance is greater than the current distance and the drone reaches the top-left corner,
		// return the bottom-right boundary coordinates.
		if maxDistance > distance && x == 1 && y == estate.Width {
			return distance, Coordinate{X: estate.Length, Y: estate.Width}
		}
		// If the distance exceeds maxDistance and the drone is at the rightmost edge,
		// move one row upward and return the new position.
		if distance > maxDistance && x == estate.Length {
			return maxDistance, Coordinate{X: x, Y: y - 1}
		}
		// If the distance exceeds maxDistance, move one step to the right and return the new position.
		if distance > maxDistance {
			return maxDistance, Coordinate{X: x + 1, Y: y}
		}
	}

	// Return zero and an empty coordinate if no conditions are met.
	return 0, Coordinate{}
}
