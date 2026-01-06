package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	// Импортируем пакет spentcalories для расчета калорий
	"github.com/mrdevelopgo/fitTracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах (для расчета дистанции в daysteps)
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// ОШИБКА ТЕСТОВ ПО ПРОБЕЛАМ ВОЗЛЕ ЗАПЯТОЙ
// БЫЛО: Парсим строку с данными о шагах и времени
//func parsePackage(data string) (int, time.Duration, error) {
// 1. Разделяем строку по запятой на части
//parts := strings.Split(data, ",")

// ПОПЫТКА: Пробуем исключить ошибку по всей строке
func parsePackage(data string) (int, time.Duration, error) {
	// 1.1. Проверяем, что вся строка не начинается пробелом
	// Тесты требуют, чтобы пробелы в начале ВСЕЙ строки вызывали ошибку
	if strings.HasPrefix(data, " ") {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
	}

	// 1.2. Проверяем, что нет пробела перед запятой
	// Например: "1234 ,1h30m" - пробел перед запятой
	if strings.Contains(data, " ,") {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
	}

	// 1.3. Проверяем, что нет пробела после запятой в начале второго значения
	// Например: "1234, 1h30m"
	if strings.Contains(data, ",") {
		parts := strings.SplitN(data, ",", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) == "" {
			return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
		}
	}

	// 1.4. Разделяем строку по запятой на части
	parts := strings.Split(data, ",")

	// 2. Проверяем, что получилось ровно 2 части (шаги и время)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
	}

	// 3. Обрабатываем первую часть - количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	// Преобразуем строку в число
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		// Используем %w для оборачивания ошибки
		return 0, 0, fmt.Errorf("неверный формат шагов: %w", err)
	}

	// 4. Проверяем, что количество шагов положительное
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}

	// 5. Обрабатываем вторую часть - время
	// Убираем лишние пробелы
	durationStr := strings.TrimSpace(parts[1])
	// Преобразуем строку в тип time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат времени: %s", durationStr)
	}

	// 6. Проверяем, что продолжительность положительная
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	// 7. Возвращаем результат: количество шагов и продолжительность
	return steps, duration, nil
}

// Обрабатываем данные о дневной активности и возвращаем информацию
func DayActionInfo(data string, weight, height float64) string {
	// 1. Парсим данные из строки
	steps, duration, err := parsePackage(data)
	if err != nil {
		// Если произошла ошибка при парсинге и логируем её
		log.Printf("Ошибка обработки данных '%s': %v", data, err)
		// Возвращаем пустую строку
		return ""
	}

	// 2. Дополнительная проверка на положительное количество шагов
	if steps <= 0 {
		return ""
	}

	// 3. Рассчитываем пройденную дистанцию
	distanceMeters := float64(steps) * stepLength
	// Переводим метры в километры
	distanceKm := distanceMeters / mInKm

	// 4. Рассчитываем количество сожженных калорий
	// По условию для дневной активности считаем как ходьбу. Я пытаюсь использовать функцию из пакета spentcalories (не уверен что я сделал это правильно)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		// Возвращем ошибку при расчете калорий
		log.Printf("Ошибка расчета калорий: %v", err)
		// Устанавливаем калории в 0
		calories = 0.0
	}

	// 5. Форматируем результат в виде строки
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)

	// 6. Возвращаем отформатированный результат
	return result
}
