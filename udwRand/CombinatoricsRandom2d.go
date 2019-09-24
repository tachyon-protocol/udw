package udwRand

import (
	"fmt"

	"github.com/tachyon-protocol/udw/udwMath"
	"github.com/tachyon-protocol/udw/udwSlice"
	"github.com/tachyon-protocol/udw/udwSort"
)

type CombinatoricsRandom2d struct {
	ANumList     []int
	BNumList     []int
	ValidCombine [][]bool
	Output       []udwMath.IntVector2
}

func (c *CombinatoricsRandom2d) Random(r *UdwRand) (err error) {
	c.Output = []udwMath.IntVector2{}
	aList := c.ANumList
	bList := c.BNumList
	validCombine := make([][]bool, len(c.ValidCombine))
	copy(validCombine, c.ValidCombine)
	sumA := 0
	sumB := 0
	var virtualType combinatoricsRandom2dVirtualType
	for _, num := range aList {
		sumA += num
	}
	for _, num := range bList {
		sumB += num
	}

	switch {
	case sumA > sumB:
		virtualType = combinatoricsRandom2dVirtualTypeB
		diff := sumA - sumB
		sumB = sumA
		bList = append(bList, diff)
		for i := range aList {
			validCombine[i] = append(validCombine[i], true)
		}
	case sumA < sumB:
		virtualType = combinatoricsRandom2dVirtualTypeA
		diff := sumB - sumA
		sumA = sumB
		aList = append(aList, diff)
		thisRow := make([]bool, len(bList))
		for i := range thisRow {
			thisRow[i] = true
		}
		validCombine = append(validCombine, thisRow)
	case sumA == sumB:
		virtualType = combinatoricsRandom2dVirtualTypeNone
	}

	aValidBCombineNumList := make([]int, len(aList))
	for AKindId := range aList {
		for _, thisValid := range validCombine[AKindId] {
			if thisValid {
				aValidBCombineNumList[AKindId]++
			}
		}
	}

	ASortOrderList := udwSlice.IntRangeSlice(len(aList))
	udwSort.IntLessCallbackSort(ASortOrderList, func(i int, j int) bool {
		i = ASortOrderList[i]
		j = ASortOrderList[j]

		switch {
		case aValidBCombineNumList[i] < aValidBCombineNumList[j]:
			return true
		case aValidBCombineNumList[i] > aValidBCombineNumList[j]:
			return false
		}

		return aList[i] > aList[j]
	})

	theOutput := []udwMath.IntVector2{}
	BRemainList := make([]int, sumB)
	BRemainListIndex := 0
	for BKindId, num := range bList {
		for i := 0; i < num; i++ {
			BRemainList[BRemainListIndex] = BKindId
			BRemainListIndex++
		}
	}

	for _, AKindId := range ASortOrderList {
		for i := 0; i < aList[AKindId]; i++ {
			BValidList := []int{}
			for _, BKindId := range BRemainList {
				if validCombine[AKindId][BKindId] {
					BValidList = append(BValidList, BKindId)
				}
			}

			if len(BValidList) == 0 {
				return fmt.Errorf("[CombinatoricsRandom2d.Random]AKindId:%d len(BValidList)==0", AKindId)
			}
			choiceBKindId := r.ChoiceFromIntSlice(BValidList)

			theOutput = append(theOutput, udwMath.IntVector2{
				X: AKindId,
				Y: choiceBKindId,
			})

			udwSlice.IntSliceRemove(&BRemainList, choiceBKindId)
		}
	}

	switch virtualType {
	case combinatoricsRandom2dVirtualTypeA:
		removeKindId := len(aList) - 1
		for _, row := range theOutput {
			if row.X != removeKindId {
				c.Output = append(c.Output, row)
			}
		}
	case combinatoricsRandom2dVirtualTypeB:
		removeKindId := len(bList) - 1
		for _, row := range theOutput {
			if row.Y != removeKindId {
				c.Output = append(c.Output, row)
			}
		}
	case combinatoricsRandom2dVirtualTypeNone:
		c.Output = theOutput
	}

	thisLen := len(c.Output)
	theOutput = make([]udwMath.IntVector2, thisLen)
	permSlice := r.Perm(thisLen)
	for i := 0; i < thisLen; i++ {
		theOutput[i] = c.Output[permSlice[i]]
	}
	c.Output = theOutput
	return
}

type combinatoricsRandom2dVirtualType int

const (
	combinatoricsRandom2dVirtualTypeA combinatoricsRandom2dVirtualType = iota
	combinatoricsRandom2dVirtualTypeB
	combinatoricsRandom2dVirtualTypeNone
)
