package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Количество метров в одном километре
	mInKm = 1000
	// Количество минут в одном часе
	minInH = 60
	// Коэффициент для расчета длины шага: длина шага = рост * 0.45
	stepLengthCoefficient = 0.45
	// Коэффициент расхода калорий при ходьбе
	walkingCaloriesCoefficient = 0.5
)

// Форматируем информацию о тренировке
func formatTrainingInfo(activity string, duration time.Duration, distance, speed, calories float64) string {
	// Переводим продолжительность в часы для вывода
	durationHours := duration.Hours()

	// Форматируем все данные в одну строку
	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationHours,
		distance,
		speed,
		calories,
	)
}

// Парсим строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделяем строку по запятым на части
	parts := strings.Split(data, ",")

	// 2. Проверяем, что получилось ровно 3 части (шаги, активность, время)
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf(
			"неправильный формат данных: ожидается 3 части, получено %d",
			len(parts),
		)
	}

	// 3. Обрабатываем первую часть - количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		// Используем %w для оборачивания ошибки
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов '%s': %w", stepsStr, err)
	}

	// 4. Проверяем, что количество шагов положительное
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// 5. Обрабатываем вторую часть - тип активности
	activity := strings.TrimSpace(parts[1])

	// 6. Обрабатываем третью часть - время тренировки
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Используем %w для оборачивания ошибки
		return 0, "", 0, fmt.Errorf("ошибка парсинга времени '%s': %w", durationStr, err)
	}

	// 7. Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// 8. Возвращаем результат
	return steps, activity, duration, nil
}

// Рассчитываем пройденную дистанцию в километрах
func distance(steps int, height float64) float64 {
	// 1. Рассчитываем длину одного шага
	stepLength := height * stepLengthCoefficient

	// 2. Рассчитываем общую дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 3. Переводим метры в километры
	distanceKilometers := distanceMeters / float64(mInKm)

	return distanceKilometers
}

// Рассчитываем среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0
	}

	// 2. Рассчитываем дистанцию
	distanceKilometers := distance(steps, height)

	// 3. Переводим продолжительность в часы
	durationHours := duration.Hours()

	// 4. Рассчитываем скорость: расстояние / время
	speed := distanceKilometers / durationHours

	return speed
}

// Рассчитываем количество сожженных калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверяем входные данные
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным: %.2f", weight)
	}

	if height <= 0.5 || height > 2.5 {
		return 0, fmt.Errorf("рост должен быть в пределах от 0.5 до 2.5 метров: %.2f", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	// 2. Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// 4. Рассчитываем калории по формуле для бега
	calories := (weight * speed * durationMinutes) / float64(minInH)

	// 5. Возвращаем результат
	return calories, nil
}

// Рассчитываем количество сожженных калорий при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверяем входные данные
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным: %.2f кг", weight)
	}

	if height <= 0.5 || height > 2.5 {
		return 0, fmt.Errorf("рост должен быть в пределах от 0.5 до 2.5 метров: %.2f м", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной: %v", duration)
	}

	// 2. Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// 4. Рассчитываем базовые калории по формуле
	baseCalories := (weight * speed * durationMinutes) / float64(minInH)

	// 5. Для ходьбы умножаем на коэффициент 0.5
	calories := baseCalories * walkingCaloriesCoefficient

	// 6. Возвращаем результат
	return calories, nil
}

// Обрабатываем данные о тренировке и возвращаем в отформатированном виде
func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Парсим данные из строки
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		// Логируем ошибку парсинга
		log.Println("Ошибка парсинга тренировки:", err)
		// Возвращаем ошибку
		return "", err
	}

	// 2. Рассчитываем дистанцию и скорость
	distanceKilometers := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитываем калории в зависимости от типа активности
	var calories float64

	// Приводим тип активности к нижнему регистру для сравнения
	switch strings.ToLower(activity) {
	case "ходьба":
		// Используем функцию для ходьбы
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "бег":
		// Используем функцию для бега
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		// Если тип активности неизвестен, возвращаем ошибку
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	// 4. Проверяем, не произошла ли ошибка при расчете калорий
	if err != nil {
		// Логируем ошибку расчета
		log.Println("Ошибка расчета калорий:", err)
		// Возвращаем ошибку
		return "", err
	}

	// 5. Форматируем информацию о тренировке
	result := formatTrainingInfo(activity, duration, distanceKilometers, speed, calories)

	// 6. Возвращаем результат
	return result, nil
}
