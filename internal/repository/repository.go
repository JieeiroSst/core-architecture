package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/JIeeiroSst/core-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Repositories struct {
	Schools        Schools
	Students       Students
	StudentLessons StudentLessons
	Courses        Courses
	Modules        Modules
	Packages       Packages
	LessonContent  LessonContent
	Offers         Offers
	PromoCodes     PromoCodes
	Orders         Orders
	Admins         Admins
	Users          Users
	Files          Files
	SurveyResults  SurveyResults
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools:        NewSchoolsRepo(db),
		Students:       NewStudentsRepo(db),
		StudentLessons: NewStudentLessonsRepo(db),
		Courses:        NewCoursesRepo(db),
		Modules:        NewModulesRepo(db),
		LessonContent:  NewLessonContentRepo(db),
		Offers:         NewOffersRepo(db),
		PromoCodes:     NewPromocodeRepo(db),
		Orders:         NewOrdersRepo(db),
		Admins:         NewAdminsRepo(db),
		Packages:       NewPackagesRepo(db),
		Users:          NewUsersRepo(db),
		Files:          NewFilesRepo(db),
		SurveyResults:  NewSurveyResultsRepo(db),
	}
}

func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
	var opts *options.FindOptions
	if pagination != nil {
		opts = &options.FindOptions{
			Skip:  pagination.GetSkip(),
			Limit: pagination.GetLimit(),
		}
	}

	return opts
}

func filterDateQueries(dateFrom, dateTo, fieldName string, filter bson.M) error {
	if dateFrom != "" && dateTo != "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return err
		}

		dateTo, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$and": []bson.M{
				{fieldName: bson.M{"$gte": dateFrom}},
				{fieldName: bson.M{"$lte": dateTo}},
			},
		})
	}

	if dateFrom != "" && dateTo == "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			fieldName: bson.M{"$gte": dateFrom},
		})
	}

	if dateFrom == "" && dateTo != "" {
		dateTo, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			fieldName: bson.M{"$lte": dateTo},
		})
	}

	return nil
}
