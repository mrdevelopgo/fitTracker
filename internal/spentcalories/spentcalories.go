package spentcalories

import (
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
	// TODO: реализовать функцию
	// Разделяем строку по запятой
	parts := strings.Split[data, ""]

	//Проверяем, что получилось три части
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf(
			"Неправильный формат данных",
			len(parts),
		)
	}

	// 3. Парсим количество шагов
    stepsStr := parts[0]
    steps, err := strconv.Atoi(stepsStr)
    if err != nil {
        return 0, "", 0, fmt.Errorf("ошибка парсинга шагов '%s': %v", stepsStr, err)
    }

	// 4. Проверяем, что шаги > 0
    if steps <= 0 {
        return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
    }
    
    // 5. Получаем вид активности (вторая часть)
    activity := parts[1]

	// 6. Парсим время (третья часть)
    durationStr := parts[2]
    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        return 0, "", 0, fmt.Errorf("ошибка парсинга времени '%s': %v", durationStr, err)
    }
    
    // 7. Проверяем, что время > 0
    if duration <= 0 {
        return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
    }

	// 8. Возвращаем результат
    return steps, activity, duration, nil

}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	// 1. Рассчитываем длину шага в метрах
	stepLength := height * stepLengthCoefficient

	// 2. Вычисляем дистанцию в метрах
	distanceMeters := (steps)float64 * lenStep

	// 3. Переводим метры в километры
	distanceKiloMetres := distanceMeters / mInKm

	// 4. Возвращаем результат
	return distanceKiloMetres
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	// 1. Проверка duration
	if duration <= 0 {
		return 0
	}
	
	// 2. Вычисляем дистанцию
	distanceKiloMetres := distance(steps, height) 

	// 3. Переводим продолжительность в часы
    durationHours := duration.Hours()

	// 4. Вычисляем среднюю скорость (км/ч)
    speed := distanceKiloMetres / durationHours
    
    // 5. Возвращаем скорость
    return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	// 1. Парсим данные из строки
    steps, activity, duration, err := parseTraining(data)
    if err != nil {
        log.Println("Ошибка парсинга тренировки:", err)
        return "", err
    }
    
    // 2. Вычисляем общие параметры (дистанция и скорость)
    distanceKiloMetres := distance(steps, height)
    speed := meanSpeed(steps, height, duration)
    
    // 3. Определяем тип тренировки и считаем калории
    var calories float64
    var errCal error
    
    switch strings.ToLower(activity) {
    case "ходьба", "ходьба ", "walking":
        calories, errCal = spentcalories.WalkingSpentCalories(steps, weight, height, duration)
    case "бег", "бег ", "running":
        calories, errCal = spentcalories.RunningSpentCalories(steps, weight, height, duration)
    default:
        return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
    }
    
    // 4. Проверяем ошибку расчета калорий
    if errCal != nil {
        log.Println("Ошибка расчета калорий:", errCal)
        return "", errCal
    }
    
    // 5. Форматируем результат
    result := formatTrainingInfo(activity, duration, distanceKiloMetres, speed, calories)
    
    // 6. Возвращаем результат
    return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	// 1. Проверяем входные параметры на корректность
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
    // Для этого нам нужно импортировать пакет daysteps
    speed := meanSpeed(steps, height, duration)
    
    // 3. Переводим продолжительность в минуты
    durationMinutes := duration.Minutes()
    
    // 4. Рассчитываем калории по формуле
    // (weight * meanSpeed * durationInMinutes) / minInH
    calories := (weight * speed * durationMinutes) / minInH
    
    // 5. Возвращаем результат
    return calories, nil	
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	// 1. Проверяем входные параметры на корректность
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
    
    // 3. Проверяем скорость для ходьбы
    if speed < 1.0 {
        return 0, fmt.Errorf("скорость %.2f км/ч слишком мала для ходьбы", speed)
    }
    
    if speed > 8.0 {
        return 0, fmt.Errorf("скорость %.2f км/ч слишком велика для ходьбы", speed)
    }
    
    // 4. Переводим продолжительность в минуты
    durationMinutes := duration.Minutes()
    
    // 5. Рассчитываем базовые калории (как для бега)
    baseCalories := (weight * speed * durationMinutes) / minInH
    
    // 6. Применяем корректирующий коэффициент для ходьбы
    calories := baseCalories * walkingCaloriesCoefficient
    
    // 7. Округляем результат до 2 знаков после запятой
    calories = math.Round(calories*100) / 100
    
    // 8. Возвращаем результат
    return calories, nil
}
