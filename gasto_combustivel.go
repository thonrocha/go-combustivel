package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {

	incluirCmd := flag.NewFlagSet("incluir", flag.ExitOnError)
	incluirKM := incluirCmd.Int("km", 0, "km")
	incluirCombustivel := incluirCmd.Float64("combustivel", 0.0, "combustivel")

	mediaCmd := flag.NewFlagSet("media", flag.ExitOnError)

	data := carregaDados()

	if len(os.Args) < 2 {
		fmt.Println("Execute novamente usando os subcomandos 'incluir' ou 'media'")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "incluir":
		incluirCmd.Parse(os.Args[2:])
		data = inserir(data, *incluirKM, *incluirCombustivel)
		gravar(data)
		fmt.Println("Novo registro incluído!")
	case "media":
		mediaCmd.Parse(os.Args[2:])
		mediaPercurso, mediaTotal := consultaMedia(data)
		fmt.Printf("Média do Percurso: %6.2f km/l \n", mediaPercurso)
		fmt.Printf("Média Total: %6.2f km/l \n", mediaTotal)
	default:
		fmt.Println("Execute novamente usando os subcomandos 'incluir' ou 'media'")
		os.Exit(1)
	}
}

func carregaDados() [][]string {
	csvfile, err := os.Open("data.csv")

	if err != nil {
		log.Fatalln("Não foi possível abrir o arquivo CSV", err)
	}

	reader := csv.NewReader(csvfile)

	data, err := reader.ReadAll()

	if err != nil {
		log.Fatalln("Não foi possível ler os dados do arquivoCSV:", err.Error())
	}
	return data
}

func inserir(data [][]string, km int, combustivel float64) [][]string {
	ultimoKmTotal, err := strconv.Atoi(data[len(data)-1][0])

	if err != nil {
		fmt.Println("Não foi possível ler a última entrada para o cálculo de KM de consumo")
		os.Exit(1)
	}

	kmIntervalo := km - ultimoKmTotal
	ultimoKmInicio, err := strconv.Atoi(data[len(data)-1][2])

	if err != nil {
		fmt.Println("Não foi possível ler a última entrada de KM para o cálculo de consumo")
		os.Exit(1)
	}

	kmInicio := ultimoKmInicio + kmIntervalo

	strUltimoTotalCombustivel := data[len(data)-1][3]
	ultimoTotalLitro, err := strconv.ParseFloat(strUltimoTotalCombustivel, 64)

	if err != nil {
		fmt.Println("Não foi possível ler a última entrada de combustível para o cálculo de consumo")
		os.Exit(1)
	}

	totalLitro := ultimoTotalLitro + combustivel

	mediaIntervalo := float64(kmIntervalo) / combustivel
	mediaTotal := float64(kmInicio) / totalLitro

	data = append(data, []string{strconv.Itoa(km),
		strconv.Itoa(kmIntervalo),
		strconv.Itoa(kmInicio),
		strconv.FormatFloat(combustivel, 'f', 8, 64),
		strconv.FormatFloat(totalLitro, 'f', 8, 64),
		strconv.FormatFloat(mediaIntervalo, 'f', 8, 64),
		strconv.FormatFloat(mediaTotal, 'f', 8, 64)})

	return data
}

func consultaMedia(data [][]string) (float64, float64) {
	strMediaPercurso := data[len(data)-1][5]
	strMediaTotal := data[len(data)-1][6]

	mediaPercurso, err := strconv.ParseFloat(strMediaPercurso, 64)

	if err != nil {
		fmt.Println("Não foi possível ler a média do percurso")
		os.Exit(1)
	}

	mediaTotal, err := strconv.ParseFloat(strMediaTotal, 64)

	if err != nil {
		fmt.Println("Não foi possível ler a média total")
		os.Exit(1)
	}

	return mediaPercurso, mediaTotal
}

func gravar(data [][]string) {
	f, err := os.Create("data.csv")
	if err != nil {
		log.Fatalf("Não foi possível criar o arquivo : %s\n", err.Error())
	}

	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Não foi possível gravar o arquivo : %s\n", e.Error())
		}
	}()

	w := csv.NewWriter(f)
	err = w.WriteAll(data)
}
