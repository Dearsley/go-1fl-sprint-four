package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: expected 2 fields, got %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps format: %w", err)
	}
	if steps <= 0 {
		return 0, 0, errors.New("steps can't be equal 0 or less")
	}

	walkDuration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid time duration format: %w", err)
	}
	if walkDuration <= 0 {
		return 0, 0, errors.New("walk duration can't be equal 0 or less")
	}

	return steps, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	distance := float32(steps) * stepLength
	distanceInKm := distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	outputMessage := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceInKm, calories)

	return outputMessage
}
