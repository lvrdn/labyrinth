// Для поиска оптимального пути применен алгоритм A*(astar)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Field struct {
	Relations     map[*Cell][]*Cell
	Length, Width int
	Start, Finish *Cell
	OpenedList    map[*Cell]int
	ClosedList    map[*Cell]struct{}
}

type Cell struct {
	Weight, I, J int
	Parent       *Cell
}

func (f *Field) GetCellNum(i, j int) int {
	return i + j + i*(f.Width-1)
}

// расчет эврестического приближения по методу манхетана
func (f *Field) GetMH(cell *Cell) int {
	return abs(f.Finish.I-cell.I) + abs(f.Finish.J-cell.J)
}

func (f *Field) FindCellInOpenedListWithMinWeight() *Cell {
	min := 0
	var foundedCell *Cell
	for cell, weight := range f.OpenedList {
		if _, ok := f.ClosedList[cell]; ok {
			continue
		}

		mh := f.GetMH(cell)

		totalWeight := weight + mh
		if min == 0 {
			min = totalWeight
			foundedCell = cell
		} else if min > totalWeight {
			min = totalWeight
			foundedCell = cell
		}
	}
	return foundedCell
}

func (f *Field) FindOptimalCell(cell *Cell) *Cell {

	cellFromOpenedList := f.FindCellInOpenedListWithMinWeight()
	if cell == nil {
		if cellFromOpenedList == nil {
			return nil
		} else {
			return cellFromOpenedList
		}
	}

	if (cellFromOpenedList.Weight + f.GetMH(cellFromOpenedList)) > (cell.Weight + f.GetMH(cell)) {
		return cell
	} else {
		return cellFromOpenedList
	}
}

func main() {

	field := &Field{
		Relations: make(map[*Cell][]*Cell),
	}

	err := scanInput(field)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if field.Start == field.Finish {
		fmt.Printf("%d %d\n.\n", field.Start.I, field.Start.J)
		return
	}

	field.OpenedList = make(map[*Cell]int)
	field.ClosedList = make(map[*Cell]struct{})

	currentCell := field.Start
	field.OpenedList[currentCell] = currentCell.Weight

LOOP:
	for {

		minWeightOfNextCell := 0
		var nextCellWithMinWeight *Cell

		for _, nextCell := range field.Relations[currentCell] {
			//проверка ячейки на то, что она уже пройдена
			if _, ok := field.ClosedList[nextCell]; ok {
				continue
			}
			//проверка на стену
			if nextCell.Weight == 0 {
				continue
			}

			if nextCell == field.Finish {
				field.Finish.Parent = currentCell
				break LOOP //найдена конечная точка
			}

			weight, ok := field.OpenedList[nextCell]
			newWeight := nextCell.Weight + field.OpenedList[currentCell]

			//определение веса пути от начальной клетки до nextCell и запись значения в OpenList
			if ok {
				if weight > newWeight {
					field.OpenedList[nextCell] = newWeight
					nextCell.Parent = currentCell
					weight = newWeight
				}
			} else {
				field.OpenedList[nextCell] = newWeight
				nextCell.Parent = currentCell
				weight = newWeight
			}

			mh := field.GetMH(nextCell)
			resultWeight := weight + mh

			//определение nextCell с минимальным весом пути, куда можно будет перейти из currentCell
			if minWeightOfNextCell == 0 {
				minWeightOfNextCell = resultWeight
				nextCellWithMinWeight = nextCell
			} else if minWeightOfNextCell > resultWeight {
				minWeightOfNextCell = resultWeight
				nextCellWithMinWeight = nextCell
			}
		}

		//currentCell пройдена
		field.ClosedList[currentCell] = struct{}{}

		//проверка и сравнение найденной клетки (если она не nil) с клеткой из OpenList, обладающей минимальным весом
		currentCell = field.FindOptimalCell(nextCellWithMinWeight)
		if currentCell == nil {
			break LOOP // нет прохода к конечной точке
		}
	}

	printOutput(field.Start, field.Finish)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func scanInput(field *Field) error {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	nums := strings.Split(input, " ")
	if len(nums) != 2 {
		return fmt.Errorf("необходимо ввести два числа - размер лабиринта")
	}

	var err error
	field.Length, err = strconv.Atoi(nums[0])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[0])
	}
	if field.Length <= 0 {
		return fmt.Errorf("длина лабиринта [%d] не может меньше или равной 0", field.Length)
	}

	field.Width, err = strconv.Atoi(nums[1])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[1])
	}
	if field.Width <= 0 {
		return fmt.Errorf("ширина лабиринта [%d] не может меньше или равной 0", field.Width)
	}

	cells := make([]*Cell, field.Length*field.Width)
	for i := range cells {
		cells[i] = new(Cell)
	}

	for i := range field.Length {
		scanner.Scan()
		input = scanner.Text()
		nums = strings.Split(input, " ")
		if len(nums) != field.Width {
			return fmt.Errorf("введенное количество значений не соответствует ранее введенным данным")
		}

		for j, weight := range nums {
			weightINT, err := strconv.Atoi(weight)
			if err != nil {
				return fmt.Errorf("значение ячейки поля {%d;%d} не является числом [%s]", i, j, weight)
			}
			if weightINT < 0 || weightINT > 9 {
				return fmt.Errorf("введенное значение [%d] ячейки поля {%d;%d} должно быть в интервале 0..9 включительно", weightINT, i, j)
			}

			num := field.GetCellNum(i, j)
			cell := cells[num]

			cell.Weight = weightINT
			cell.I = i
			cell.J = j

			if i != 0 {
				num := field.GetCellNum(i, j) - field.Width
				field.Relations[cells[num]] = append(field.Relations[cells[num]], cell)
			}

			if i != field.Length-1 {
				num := field.GetCellNum(i, j) + field.Width
				field.Relations[cells[num]] = append(field.Relations[cells[num]], cell)
			}

			if j != 0 {
				num := field.GetCellNum(i, j) - 1
				field.Relations[cells[num]] = append(field.Relations[cells[num]], cell)
			}

			if j != field.Width-1 {
				num := field.GetCellNum(i, j) + 1
				field.Relations[cells[num]] = append(field.Relations[cells[num]], cell)
			}
		}

	}

	scanner.Scan()
	input = scanner.Text()
	nums = strings.Split(input, " ")
	if len(nums) != 4 {
		return fmt.Errorf("необходимо ввести четыре числа - координаты начальной и конечной точек")
	}

	iStart, err := strconv.Atoi(nums[0])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[0])
	}

	jStart, err := strconv.Atoi(nums[1])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[1])
	}

	if iStart >= field.Length || iStart < 0 || jStart >= field.Width || jStart < 0 {
		return fmt.Errorf("начальная точка {%d;%d} должна иметь координаты {0..%d;0..%d} включительно", iStart, jStart, field.Length-1, field.Width-1)
	}

	num := field.GetCellNum(iStart, jStart)

	field.Start = cells[num]
	if field.Start.Weight == 0 {
		return fmt.Errorf("начальная точка {%d;%d} не может быть стеной", field.Start.I, field.Start.J)
	}

	iFinish, err := strconv.Atoi(nums[2])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[2])
	}

	jFinish, err := strconv.Atoi(nums[3])
	if err != nil {
		return fmt.Errorf("введенное значение [%s] должно быть числом", nums[3])
	}

	if iFinish >= field.Length || iFinish < 0 || jFinish >= field.Width || jFinish < 0 {
		return fmt.Errorf("конечная точка {%d;%d} должна иметь координаты {0..%d;0..%d} включительно", iFinish, jFinish, field.Length-1, field.Width-1)
	}

	num = field.GetCellNum(iFinish, jFinish)

	field.Finish = cells[num]
	if field.Start.Weight == 0 {
		return fmt.Errorf("начальная точка {%d;%d} не может быть стеной", field.Finish.I, field.Finish.J)
	}

	return nil
}

func printOutput(startCell, finishCell *Cell) {

	currentCell := finishCell
	if currentCell.Parent == nil {
		fmt.Println("Конечная точка отделена от начальной непроходимой стеной, найти путь невозможно")
		return
	}

	var output string

	for {
		output = fmt.Sprintf("%d %d\n", currentCell.I, currentCell.J) + output
		if currentCell == startCell {
			output += ".\n"
			break
		}
		currentCell = currentCell.Parent
	}

	fmt.Println(output)
}
