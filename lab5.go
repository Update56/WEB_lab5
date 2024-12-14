package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Структура для GPU
type GPU struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	MemorySize       int    `json:"memory_size"`
	CoreClock        string `json:"core_clock"`
	CudaCores        int    `json:"cuda_cores"`
	PowerConsumption int    `json:"power_consumption"`
	LinkImage        string `json:"link_image"`
}

// Структура для CPU
type CPU struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Cores     int    `json:"cores"`
	Frequency string `json:"frequency"`
	CacheSize string `json:"cache_size"`
	TDP       int    `json:"tdp"`
	LinkImage string `json:"link_image"`
}

// Функция для получения и парсинга JSON по URL
func getJSON(url string, target interface{}) error {
	// Отправка запроса
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Парсинг JSON в структуру
	return json.Unmarshal(body, target)
}

// Изменение функции для CPU
func createCPULayout(cpus []CPU, url string) *fyne.Container {
	cpuList := container.NewVBox()
	for _, cpu := range cpus {
		// Парсинг изображения
		imageURI, err := storage.ParseURI(fmt.Sprintf("%s%s", url, cpu.LinkImage))
		if err != nil {
			log.Printf("Failed to parse image URI: %v", err)
			continue
		}
		cpuImage := canvas.NewImageFromURI(imageURI)
		cpuImage.SetMinSize(fyne.NewSize(150, 150)) // Установка минимального размера изображения
		cpuImage.FillMode = canvas.ImageFillContain // Сохранение пропорций изображения

		// Создание карточки с изображением и информацией
		cpuItem := widget.NewCard(
			cpu.Name,
			fmt.Sprintf("Cores: %d | Frequency: %s | Cache: %s | TDP: %dW", cpu.Cores, cpu.Frequency, cpu.CacheSize, cpu.TDP),
			cpuImage,
		)
		cpuList.Add(cpuItem)
	}
	return cpuList
}

// Изменение функции для GPU
func createGPULayout(gpus []GPU, url string) *fyne.Container {
	gpuList := container.NewVBox()
	for _, gpu := range gpus {
		imageURI, err := storage.ParseURI(fmt.Sprintf("%s%s", url, gpu.LinkImage))
		if err != nil {
			log.Printf("Failed to parse image URI: %v", err)
			continue
		}
		gpuImage := canvas.NewImageFromURI(imageURI)
		gpuImage.SetMinSize(fyne.NewSize(150, 150)) // Установка минимального размера изображения
		gpuImage.FillMode = canvas.ImageFillContain // Сохранение пропорций изображения

		// Создание карточки с изображением и информацией
		gpuItem := widget.NewCard(
			gpu.Name,
			fmt.Sprintf("Memory: %dGB | Clock: %s | CUDA Cores: %d | TDP: %dW", gpu.MemorySize, gpu.CoreClock, gpu.CudaCores, gpu.PowerConsumption),
			gpuImage,
		)
		gpuList.Add(gpuItem)
	}
	return gpuList
}

func main() {
	var ipAddress string
	// Инициализация приложения
	myApp := app.New()
	myWindow := myApp.NewWindow("notDNSshop")
	startWin := myApp.NewWindow("nds")

	inputIPLabel := widget.NewLabel("Введите IP:port сервера")
	IPText := widget.NewEntry()
	//кнопка ввода данных
	inputBtn := widget.NewButton("Ввод", func() {
		ipAddress = IPText.Text
		Load(myWindow, ipAddress)
		myWindow.Show()
		startWin.Close()
	})
	//группировка виджетов в один контейнер
	inputCont := container.NewGridWithRows(3, inputIPLabel, IPText, inputBtn)
	//установка контенка для стартового окна
	startWin.SetContent(inputCont)
	//отображаем окно
	startWin.Show()

	myApp.Run()
}

func Load(myWindow fyne.Window, ipAddress string) {
	// URL для получения данных о GPU и CPU

	gpuURL := fmt.Sprintf("http://%s:3000/gpus", ipAddress)
	cpuURL := fmt.Sprintf("http://%s:3000/cpus", ipAddress)

	// Слайсы для хранения данных
	var gpus []GPU
	var cpus []CPU

	// Получение и парсинг данных о GPU
	err := getJSON(gpuURL, &gpus)
	if err != nil {
		log.Fatalf("Error fetching GPUs: %v", err)
	}

	// Получение и парсинг данных о CPU
	err = getJSON(cpuURL, &cpus)
	if err != nil {
		log.Fatalf("Error fetching CPUs: %v", err)
	}

	// Главная страница
	homeLabel := widget.NewLabel("Добро пожаловать в наш магазин ПК-комплектующих notDNSshop\nВыберите категорию на панели")

	// Страница с CPU
	cpuPage := createCPULayout(cpus, cpuURL)

	// Страница с GPU
	gpuPage := createGPULayout(gpus, gpuURL)

	var menu *fyne.Container

	buttons := []fyne.Widget{
		widget.NewButton("Главная", func() {
			myWindow.SetContent(container.NewBorder(menu, nil, nil, nil, homeLabel))
		}),
		widget.NewButton("CPU", func() {
			myWindow.SetContent(container.NewBorder(menu, nil, nil, nil, cpuPage))
		}),
		widget.NewButton("GPU", func() {
			myWindow.SetContent(container.NewBorder(menu, nil, nil, nil, gpuPage))
		})}

	for _, item := range buttons {
		item.Resize(fyne.NewSize(50, 100))
	}
	// Создание бокового меню
	menu = container.NewHBox(
		buttons[0], buttons[1], buttons[2],
	)
	// Установка начальной страницы
	myWindow.SetContent(container.NewBorder(menu, nil, nil, nil, homeLabel))
}
