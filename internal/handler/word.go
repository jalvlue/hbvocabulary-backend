package handler

import (
	"HBVocabulary/common"
	"HBVocabulary/internal/model"
	"HBVocabulary/token"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	roundOneWordSize = 40
)

func (server *Server) getWordListRoundOne(ctx *gin.Context) {
	totalWords := 8000
	numLayers := 10
	weights := []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1, 0.05, 0.03, 0.02, 0.01}

	layerSizes := make([]int, numLayers)
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(totalWords) * weights[i] / totalWeight)
	}

	randCount := 0
	shouldRand := make([]int, numLayers)
	for i := 0; i < numLayers; i++ {
		tmpPer := weights[i] / totalWeight
		tmpSize := tmpPer * float64(roundOneWordSize)
		shouldRand[i] = int(tmpSize)
		randCount += shouldRand[i]
	}

	if randCount < roundOneWordSize {
		shouldRand[0] += roundOneWordSize - randCount
	} else if len(shouldRand) > roundOneWordSize {
		shouldRand[0] -= randCount - roundOneWordSize
	}

	randomNumbers2 := make([]int, 0)
	curCount := 0
	for i := 0; i < numLayers; i++ {
		curCount += layerSizes[i]
		randomNumbers := generateRandomNumbers(shouldRand[i], curCount-layerSizes[i]+1, curCount)
		randomNumbers2 = append(randomNumbers2, randomNumbers...)
	}

	// 对 randomNumbers2 进行排序
	sort.Ints(randomNumbers2)

	list, err := server.store.GetWordListById(randomNumbers2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "success",
		"data": list,
	})
}

func generateRandomNumbers(n, min, max int) []int {
	if max-min+1 < n {
		log.Println("could not generate enough non-repeated numbers")
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	numbers := make([]int, n)
	used := make(map[int]bool)

	for i := 0; i < n; i++ {
		randomNum := rand.Intn(max-min+1) + min
		for used[randomNum] {
			randomNum = rand.Intn(max-min+1) + min
		}
		numbers[i] = randomNum
		used[randomNum] = true
	}

	return numbers
}

type requestWordItem struct {
	ID    int    `json:"id"`
	Word  string `json:"word"`
	Known int    `json:"known"`
}

type requestWordList struct {
	WordList []requestWordItem `json:"wordList"`
}

func (server *Server) getWordListRoundTwo(ctx *gin.Context) {
	var roundOneMsg requestWordList
	if err := ctx.ShouldBindJSON(&roundOneMsg); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(http.ErrAbortHandler))
		return
	}

	numLayers := 6
	totalWords := roundOneMsg.WordList[len(roundOneMsg.WordList)-1].ID

	weights := []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1}
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}

	knownWeights := []float64{1.0, 1.3, 1.6, 1.9, 2.2, 2.5}
	totalKnownWeights := 0.0
	for i := 0; i < numLayers; i++ {
		totalKnownWeights += knownWeights[i]
	}

	boundaries := make([]int, numLayers+1)
	boundaries[0] = 1
	for i := 1; i <= numLayers; i++ {
		boundary := int(math.Round(float64(totalWords) * (weights[i-1] / totalWeight)))
		boundaries[i] = boundaries[i-1] + boundary
	}

	knownPers := make([]float64, numLayers)
	totalPer := 0.0
	curID := 0
	length := len(roundOneMsg.WordList)

	for i := 0; i < numLayers; i++ {
		tmp := curID
		knownCount := 0
		layerBoundaryStart := boundaries[i]
		layerBoundaryEnd := boundaries[i+1]

		for ; curID < length && roundOneMsg.WordList[curID].ID < layerBoundaryEnd; curID++ {
			if roundOneMsg.WordList[curID].ID > layerBoundaryStart && roundOneMsg.WordList[curID].Known == 1 {
				knownCount += 1
			}
		}

		if curID > 0 {
			knownPer := float64(knownCount) / float64(curID-tmp)
			knownPers[i] = knownPer
			totalPer += knownPer
		} else {
			knownPers[i] = 0.0
		}
	}

	totalPer = totalPer / float64(numLayers)
	sampleSize := 0
	wordsRange := 5000

	if totalPer < 0 || totalPer > 1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid know ratio",
		})
	}

	if totalPer < 0.3 {
		sampleSize = 30
		wordsRange = 5000
		weights = []float64{1.0, 0.8, 0.6, 0.4, 0.2, 0.1}
	} else if totalPer < 0.5 {
		sampleSize = 60
		wordsRange = 15000
		weights = []float64{1.0, 0.8, 0.7, 0.6, 0.5, 0.4}
	} else if totalPer < 0.8 {
		sampleSize = 80
		wordsRange = 25000
		weights = []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5}
	} else {
		sampleSize = 100
		wordsRange = 41000
		weights = []float64{1.0, 0.9, 0.9, 0.9, 0.8, 0.7}
	}

	layerSizes := make([]int, numLayers)
	totalWeight2 := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight2 += weights[i]
	}
	for i := 0; i < numLayers; i++ {
		layerSizes[i] = int(float64(wordsRange) * weights[i] / totalWeight2)
	}

	randCount := 0
	shouldRand := make([]int, numLayers)
	for i := 0; i < numLayers; i++ {
		tmpPer := weights[i] / totalWeight2
		tmpSize := tmpPer * float64(sampleSize)
		shouldRand[i] = int(tmpSize)
		randCount += shouldRand[i]
	}

	if randCount < sampleSize {
		shouldRand[0] += sampleSize - randCount
	} else if len(shouldRand) > sampleSize {
		shouldRand[0] -= randCount - sampleSize
	}

	rand.Seed(time.Now().UnixNano())

	randomNumbers2 := make([]int, 0)
	curCount := 0
	for i := 0; i < numLayers; i++ {
		curCount += layerSizes[i]
		randomNumbers := generateRandomNumbers(shouldRand[i], curCount-layerSizes[i]+1, curCount)
		randomNumbers2 = append(randomNumbers2, randomNumbers...)
	}

	sort.Ints(randomNumbers2)

	list, err := server.store.GetWordListById(randomNumbers2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "success",
		"data": list,
	})
}

type PostWords struct {
	model.Vocabulary
	Known string `json:"known"`
}

func (server *Server) getResult(ctx *gin.Context) {
	var roundTwoMsg requestWordList
	if err := ctx.ShouldBindJSON(&roundTwoMsg); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	list := make([]PostWords, len(roundTwoMsg.WordList))
	for i := 0; i < len(roundTwoMsg.WordList); i++ {
		list[i].ID = uint(roundTwoMsg.WordList[i].ID)
		list[i].Word = roundTwoMsg.WordList[i].Word
		list[i].Known = strconv.Itoa(roundTwoMsg.WordList[i].Known)
	}

	sortByID(list)

	testResult := estimateVocabulary(list)

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(payload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user.TestCount += 1
	if testResult > user.MaxScore {
		user.MaxScore = testResult
	}
	err = server.store.SetTestResult(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "success",
		"data": testResult,
	})
}

func sortByID(wordList []PostWords) {
	less := func(i, j int) bool {
		return wordList[i].ID < wordList[j].ID
	}

	sort.Slice(wordList, less)
}

func estimateVocabulary(list []PostWords) int {
	// 设置抽样参数
	numLayers := 10 // 词汇太少了，无法估算
	if len(list) < numLayers {
		return 20
	}
	//totalWords := 54000
	totalWords := int(list[len(list)-1].ID)

	// 定义层次权重，采用指数衰减的方式
	weights := []float64{1.0, 0.9, 0.8, 0.7, 0.6, 0.5, 0.4, 0.3, 0.1, 0.1}
	totalWeight := 0.0
	for i := 0; i < numLayers; i++ {
		totalWeight += weights[i]
	}

	// 定义层次的得分权重
	knownWeights := []float64{2.0, 2.1, 2.2, 2.3, 2.35, 2.4, 2.45, 2.5, 2.6, 2.7}

	// 定义层次系数
	//levelCoefficient := []float64{0.95, 0.8, 0.85, 0.8, 0.75, 0.7}

	// 计算每个层次的分界线 = 划分层级
	boundaries := make([]int, numLayers+1)
	boundaries[0] = 1
	for i := 1; i <= numLayers; i++ {
		boundary := int(math.Round(float64(totalWords) * (weights[i-1] / totalWeight)))
		boundaries[i] = boundaries[i-1] + boundary
	}

	// 计算每个层次有多少单词认识 = 认识率
	knownPers := make([]float64, numLayers)
	totalPer := 0.0
	curId := 0
	length := len(list)

	for i := 0; i < numLayers; i++ {
		tmp := curId
		knownCount := 0
		layerBoundaryStart := boundaries[i]
		layerBoundaryEnd := boundaries[i+1]

		for ; curId < length && list[curId].ID < uint(layerBoundaryEnd); curId++ {
			if list[curId].ID > uint(layerBoundaryStart) && list[curId].Known == "1" {
				knownCount++
			}
		}

		// 计算该层级的认识率
		if curId > 0 {
			knownPer := float64(knownCount) / float64(curId-tmp)
			knownPers[i] = knownPer
		} else {
			knownPers[i] = 0.0
		}
	}

	// 避免某些中间层级没有单词，导致平均认识率加了NAN，为0，这里做一下处理
	skip := 0
	for _, num := range knownPers {
		if math.IsNaN(num) {
			skip++
			continue // 跳过NaN值
		}
		totalPer += num
	}

	totalPer = totalPer / float64(numLayers-skip)
	if totalPer == 0 {
		return 20
	}

	totalKnownWeights := 0.0
	for i := 0; i < len(knownPers); i++ {
		if math.IsNaN(knownPers[i]) {
			skip++
			continue // 跳过NaN值
		}
		totalKnownWeights += knownWeights[i]
	}

	// 根据认识率的得分权重计算每个层级的得分
	knownScores := make([]float64, numLayers)
	totalScore := 0.0
	for i := 0; i < numLayers; i++ {
		if math.IsNaN(knownPers[i]) {
			continue // 跳过NaN值
		}
		knownScores[i] = knownPers[i] * (knownWeights[i] / totalKnownWeights)
		totalScore += knownScores[i]
	}
	totalNum2 := math.Round(float64(totalWords) * totalScore)
	// float64 转换成 int
	totalNum := int(totalNum2)
	if totalNum > 54000 || totalNum < 0 {
		return 20
	}

	return totalNum

}
