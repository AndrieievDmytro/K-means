package main

import (
	"fmt"
	"math"
	"os"
	"sort"
)

var (
	clusterNum = 3
	trainFile  = "./data/train.txt"
)

func (fl *Flowers) sortStruct() {
	sort.SliceStable(fl.Fl, func(i, j int) bool {
		for idx := range fl.Fl[i].Params {
			if fl.Fl[i].Params[idx] < fl.Fl[j].Params[idx] {
				return true
			} else if fl.Fl[i].Params[idx] > fl.Fl[j].Params[idx] {
				return false
			}
		}
		return false
	})
}

func (fl *Flowers) AssignRandomObservations() map[int]map[int][]float64 {
	clusters := make(map[int]map[int][]float64)

	for i := 0; i < clusterNum; i++ {
		clusters[i] = make(map[int][]float64)
	}
	clNum := 0
	clusterElementsNum := len(fl.Fl) / (clusterNum)
	for ind, val := range fl.Fl {
		if ind > 0 && ind%clusterElementsNum == 0 {
			if clNum < (clusterNum - 1) {
				clNum++
			}
		}
		clusters[clNum][ind-clNum*clusterElementsNum] = make([]float64, len(val.Params))
		copy(clusters[clNum][ind-clNum*clusterElementsNum], val.Params)
	}
	return clusters
}

// as parameter map[clusterIndex]map[observation line from data set][array of data sample parametres]float64
// as a return value map[clasterIdx]distance from observation to the centroid
func CalculateCentroids(clusters map[int]map[int][]float64) map[int][]float64 {

	centroids := make(map[int][]float64, clusterNum)

	for clusterInd, observations := range clusters {

		for _, observation := range observations {
			if _, ok := centroids[clusterInd]; !ok {
				centroids[clusterInd] = make([]float64, len(observation))
				copy(centroids[clusterInd], observation)
			} else {
				for i := range observation {
					centroids[clusterInd][i] += observation[i]
				}
			}
		}
	}

	for ind, array := range centroids {
		denominator := float64(len(clusters[ind]))
		for indx := range array {
			centroids[ind][indx] = array[indx] / denominator
		}
	}
	return centroids
}

func CalculateDistances(observation []float64, centroid []float64) float64 {
	if len(observation) != len(centroid) {
		fmt.Println("Different dimentions")
		os.Exit(1)
	}
	var distance float64
	for i, val := range observation {
		distance += math.Pow(val-centroid[i], 2)
	}
	return distance
}

func GetMinVal(values []map[int]float64) (int, float64) {
	minValInd := 0
	minVal := values[0][0]
	for _, val := range values {
		for ind := range val {
			if minVal > val[ind] {
				minVal = val[ind]
				minValInd = ind
			}
		}

	}
	return minValInd, minVal
}

func IsEqualMap(clusters map[int]map[int][]float64, newCluster map[int]map[int][]float64) bool {

	for c, val := range clusters {
		if len(val) != len(newCluster[c]) {
			return false
		}
		for k := range val {
			if _, ok := newCluster[c][k]; !ok {
				return false
			}
		}
	}
	return true
}

func PrintStat(tr *Flowers, clusters map[int]map[int][]float64) {

	for c, cv := range clusters {
		fmt.Println(c, "(len:", len(cv), ")", ":")
		stat := make(map[string]float64)
		for key := range cv {
			stat[tr.Fl[key].Name] += 1
		}
		commonCount := 0.0
		for _, val := range stat {
			commonCount += val
		}
		for key, val := range stat {
			fmt.Println("\tStat:", key, ":", val*100.0/commonCount, "%", ", Number of flowers", val)
		}
	}
	fmt.Println("-------------------------------------------")
}

func GetResults() {
	newCluster := make(map[int]map[int][]float64)
	tr := new(Flowers)
	tr.readData(trainFile) // write data to a structure correct
	tr.sortStruct()

	clusters := tr.AssignRandomObservations() // random assigned observations to clussters

	centroids := CalculateCentroids(clusters) // Initial values of centroids

	cnt := 0

	for {

		//Flowers percentage
		cnt++
		fmt.Println("Iterations:", cnt)
		PrintStat(tr, clusters)

		results := make(map[int][]map[int]float64, len(tr.Fl)) // map[line from dataSet][array of clusters]map[claster][result]

		if IsEqualMap(clusters, newCluster) {
			break
		}

		if len(newCluster) > 0 {
			clusters = make(map[int]map[int][]float64)
			for key, val := range newCluster {
				clusters[key] = val
			}
		}
		newCluster = make(map[int]map[int][]float64)

		for ind, sample := range tr.Fl {
			obsResult := make(map[int]float64)
			for centroidInd, centroidCopy := range centroids {
				obsResult[centroidInd] = CalculateDistances(sample.Params, centroidCopy)
			}
			results[ind] = append(results[ind], obsResult)
		}

		for i := 0; i < clusterNum; i++ {
			newCluster[i] = make(map[int][]float64)
		}

		for lineIndx, val := range results {
			ind, _ := GetMinVal(val)
			newCluster[ind][lineIndx] = make([]float64, len(tr.Fl[lineIndx].Params))
			copy(newCluster[ind][lineIndx], tr.Fl[lineIndx].Params)
		}
		centroids = CalculateCentroids(newCluster)
	}
}
