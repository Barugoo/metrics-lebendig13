package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	models "github.com/lebendig13/metrics/internal/model"
)

type Server struct {
	memStorage *models.MemStorage
}

func NewServer(memStorage *models.MemStorage) *Server {
	return &Server{
		memStorage: memStorage,
	}
}

func (s *Server) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := req.URL.Path[1:]
	pathSegments := strings.Split(path, "/")
	pathSegmentsLen := len(pathSegments)

	if pathSegmentsLen < 2 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad request: no metric type")
		return
	}

	mType := pathSegments[1]
	if mType != models.Counter && mType != models.Gauge {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad request: bad type")
		return
	}
	fmt.Println("mType: ", mType)

	if pathSegmentsLen < 3 {
		res.WriteHeader(http.StatusNotFound)
		fmt.Println("Bad request: no metric type")
		return
	}

	mName := pathSegments[2]
	// TODO: Проверять имя метрики
	// ...
	if mName == "" {
		res.WriteHeader(http.StatusNotFound)
		fmt.Println("Bad request: empty metric name")
		return
	}
	fmt.Println("mName: ", mName)

	if pathSegmentsLen != 4 {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad request: pathSegmentsLen = ", pathSegmentsLen)
		return
	}

	mValue := pathSegments[3]
	if mValue == "" {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad request: empty value")
		return
	}
	fmt.Println("mValue: ", mValue)

	var metric models.Metrics
	metric.ID = mName
	metric.MType = mType

	switch mType {
	case models.Counter:
		value, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Println("Bad request: cannot parse counter value")
			return
		}
		metric.Delta = &value
	case models.Gauge:
		value, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Println("Bad request: cannot parse gauge value")
			return
		}
		metric.Value = &value
	}

	err := s.memStorage.InsertOrUpdate(metric)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Internal server error: cannot save metric to storage")
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
}
