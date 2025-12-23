package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	//parts разделяет слайс на строки по запятой
	parts := strings.Split(data, ",")

	//Проверка количества частей равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
	}

	//Преобразование количества шагов в тип int
	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат шагов: %s", stepsStr)
	}

	//Проверка, что количество шагов > 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}

	//Преобразование времени (вторая часть) в time.Duration
	durationStr := parts[1]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат времени: %s", durationStr)
	}

	//Проверка, что продолжительность > 0
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	//Возвращение успешного результата
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	// 1. Получаем данные через parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
		// Выводим ошибку и возвращаем пустую строку
		fmt.Println("Ошибка:", err)
		return ""
	}

	// 2. Проверяем, что шаги > 0 (двойная проверка)
	if steps <= 0 {
		return ""
	}

	// 3. Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4. Переводим в километры
	distanceKm := distanceMeters / mInKm

	// 5. Вычисляем калории
	// TODO: Здесь будет вызов функции из spentcalories
	calories := 0.0 // временно

	// 6. Формируем строку результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)

	return result

}
