package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid format: expected 3 fields, got %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps format: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps can't be equal 0 or less")
	}

	walkDuration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}
	if walkDuration <= 0 {
		return 0, "", 0, errors.New("walk duration can't be equal 0 or less")
	}

	trainingType := parts[1]

	return steps, trainingType, walkDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceInKm := (stepLength * float64(steps)) / mInKm

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	meanSpeed := distance / duration.Hours()

	return meanSpeed
}

func validateParameters(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("steps can't be equal 0 or less")
	}
	if weight <= 0 {
		return errors.New("weight can't be equal 0 or less")
	}
	if height <= 0 {
		return errors.New("height can't be equal 0 or less")
	}
	if duration <= 0 {
		return errors.New("duration can't be equal 0 or less")
	}

	return nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	switch trainingType {
	case "Бег":
		distance := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		spentCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("failed calculate running spent calories: %w", err)
		}
		outputMessage := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, duration.Hours(), distance, speed, spentCalories)

		return outputMessage, nil

	case "Ходьба":
		distance := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		spentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("failed calculate walking spent calories: %w", err)
		}
		outputMessage := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, duration.Hours(), distance, speed, spentCalories)

		return outputMessage, nil

	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParameters(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := meanSpeed(steps, height, duration)
	spentCalories := (weight * meanSpeed * duration.Minutes()) / minInH

	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParameters(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := meanSpeed(steps, height, duration)
	spentCalories := (weight * meanSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient

	return spentCalories, nil
}
