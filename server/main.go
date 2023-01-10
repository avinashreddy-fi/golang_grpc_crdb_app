package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pb "moviesapp.com/grpc/protos"

	"google.golang.org/grpc"

	"fmt"
)

const (
	port = ":50051"
)

// used array here
var movies []*pb.MovieInfo

// db schema
type MovieInfoDb struct {
	Id        string
	Isbn      string
	Title     string
	Firstname string
	Lastname  string
}

type movieServer struct {
	pb.UnimplementedMovieServer
}

var db *gorm.DB

func main() {
	// db setup
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")+"&application_name=$ docs_simplecrud_gorm"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DB Connected!! at '%s'", time.Now())

	// Automatically creates the "MovieInfoDb" table based on the `MovieInfoDb`
	// model.
	db.AutoMigrate(&MovieInfoDb{})
	if err := initMovies(); err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterMovieServer(s, &movieServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initMovies() error {
	//movie1 := &pb.MovieInfo{Id: "1", Isbn: "0593310438",
	//	Title: "The Batman", Director: &pb.Director{
	//		Firstname: "Matt", Lastname: "Reeves"}}
	//movie2 := &pb.MovieInfo{Id: "2", Isbn: "3430220302",
	//	Title: "Doctor Strange in the Multiverse of Madness",
	//	Director: &pb.Director{Firstname: "Sam",
	//		Lastname: "Raimi"}}
	//
	//movies = append(movies, movie1)
	//movies = append(movies, movie2)

	// DB insert 2 fields
	log.Printf("Creating 2 movies info")
	//length := len(movies)
	if err := db.Create(&MovieInfoDb{
		Id:        "2",
		Isbn:      "0593310438",
		Title:     "The",
		Firstname: "Matt",
		Lastname:  "Reeves",
	}).Error; err != nil {
		return err
	}
	return nil
}

func (s *movieServer) GetMovies(in *pb.Empty,
	stream pb.Movie_GetMoviesServer) error {
	log.Printf("Received: %v", in)
	for _, movie := range movies {
		if err := stream.Send(movie); err != nil {
			return err
		}
	}
	return nil
}

func (s *movieServer) GetMovie(ctx context.Context,
	in *pb.Id) (*pb.MovieInfo, error) {
	log.Printf("Received: %v", in)

	//res := &pb.MovieInfo{}
	//
	//for _, movie := range movies {
	//	if movie.GetId() == in.GetValue() {
	//		res = movie
	//		break
	//	}
	//}
	//
	//return res, nil
	// DB fetch
	var movieInfos []MovieInfoDb
	db.Find(&movieInfos)

	res := &pb.MovieInfo{}
	fmt.Printf("Get Movie call at '%s':\n", time.Now())
	for _, movieInfo := range movieInfos {
		if movieInfo.Id == in.GetValue() {
			res.Id = movieInfo.Id
			res.Isbn = movieInfo.Isbn
			res.Title = movieInfo.Title
			res.Director = &pb.Director{Firstname: movieInfo.Firstname, Lastname: movieInfo.Lastname}
			break
		}
	}
	return res, nil
}

func (s *movieServer) CreateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Id, error) {
	log.Printf("Received: %v", in)
	res := pb.Id{}
	res.Value = strconv.Itoa(rand.Intn(100000000))
	in.Id = res.GetValue()
	movies = append(movies, in)

	if err := db.Create(&MovieInfoDb{
		Id:        string(res.GetValue()),
		Isbn:      in.GetIsbn(),
		Title:     in.GetTitle(),
		Firstname: in.GetDirector().GetFirstname(),
		Lastname:  in.GetDirector().GetLastname(),
	}).Error; err != nil {
		return &res, err
	}
	return &res, nil
}

func (s *movieServer) UpdateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}
	//for index, movie := range movies {
	//	if movie.GetId() == in.GetId() {
	//		movies = append(movies[:index], movies[index+1:]...)
	//		in.Id = movie.GetId()
	//		movies = append(movies, in)
	//		res.Value = 1
	//		break
	//	}
	//}
	var updateMovie MovieInfoDb
	db.First(&updateMovie, in.GetId())

	updateMovie.Title = in.GetTitle()
	updateMovie.Isbn = in.GetIsbn()

	if err := db.Save(&updateMovie).Error; err != nil {
		return &res, err
	}
	res.Value = 1
	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}
	//for index, movie := range movies {
	//	if movie.GetId() == in.GetValue() {
	//		movies = append(movies[:index], movies[index+1:]...)
	//		res.Value = 1
	//		break
	//	}
	//}

	err := db.Where("Id = ?", in.GetValue()).Delete(&MovieInfoDb{}).Error
	if err != nil {
		return &res, err
	}
	res.Value = 1
	return &res, nil
}
