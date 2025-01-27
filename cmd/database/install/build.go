package database

// bootstrap application data on a local Couchbase cluster
// datamodel Version 0.1

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
	"bitbought.net/project-root-directory/internal/hashcode"
)

var connectionString string = "localhost"

type DataModel struct {
	_Bucket      string
	_Scopes      []string
	_Collections []Collect
}

type Collect struct {
	CollectionName string
	ScopeName      string
}

	// define a store document

	type Geolocation struct {
		Lat float32 `json:"lat"` //7 decimal precision
		Lon float32 `json:"lon"` //7 decimal precision
	}

	type Hours struct {
		Open  [2]int `json:"open"`  // {9,30}
		Close [2]int `json:"close"` // {22,0}
	}

	type Item struct {
		Name        string  `json:"name"`
		Price       float32 `json:"price"`
		Stock       uint8   `json:"stock"`
		Description string  `json:"description"`
	}

	type BankAccount struct {
		Name    string `json:"name"`
		Account string `json:"account"`
		Transit string `json:"transit"`
		Inst    string `json:"inst"`
	}

	type Store struct {
		Version    string      `json:"_version"`
		Type       string      `json:"_type"`
		Created    int64       `json:"_created"` // epoch time in seconds
		CreatedBy  string      `json:"_createdBy"`
		Modified   int64       `json:"_modified"` // epoch time in seconds
		ModifiedBy string      `json:"_modifiedBy"`
		Store_id  string      `json:"store_id"`
		Name       string      `json:"name"`
		Phone      string      `json:"phone"`
		Address    string      `json:"address"`
		Email      string      `json:"email"`
		Fax        string      `json:"fax"`
		Url        string      `json:"url"`
		Geo        Geolocation `json:"geo"`
		City       string      `json:"city"`
		Prov       string      `json:"prov"`
		Country    string      `json:"country"`
		Hour       Hours       `json:"hour"`
		Listing    []Item      `json:"listing"`
		Bank       BankAccount `json:"bank"`
	}

func main() {

	// For a secure cluster connection, use `couchbases://<your-cluster-ip>` instead.
	cluster, err := gocb.Connect("couchbase://"+connectionString, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: "class_a_rw",
			Password: "123456",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// define data model
	dataModel := DataModel{
		_Bucket: "market",
		_Scopes: []string{"persona", "activity", "metrics", "counter"},
		_Collections: []Collect{
			{"shopper", "persona"},
			{"store", "persona"},
			{"admin", "persona"},
			{"authenticate_shopper", "activity"},
			{"authenticate_store", "activity"},
			{"authenticate_admin", "activity"},
			{"transact", "activity"},
			{"support", "activity"},
			{"worklog", "metrics"},
			{"shopper_counter", "counter"},
			{"store_counter", "counter"},
			{"admin_counter", "counter"},
			{"authenticate_shopper_counter", "counter"},
			{"authenticate_store_counter", "counter"},
			{"authenticate_admin_counter", "counter"},
			{"transact_counter", "counter"},
			{"support_counter", "counter"},
			{"worklog_counter", "counter"},
		},
	}
	// setup database in cluster
	cluster, err = setup_database(dataModel, cluster)
	if err != nil {
		panic(err)
	}
	fmt.Println("data model built success on couchbase cluster at " + time.Now().String())
    fmt.Println("loading sample data...")
	// populate database with sample data


	var id string = counter("store", cluster)
	var traderJoe Store = Store{
		Version:   "0.1.0",
		Type:      "store",
		Created:   time.Now().Unix(),
		CreatedBy: "class_a_rw",
		Store_id: id,
		Name:      "tim hortons",
		Phone:     "2361113333",
		Address:   "4949 canada way",
		Email:     "thliquor@shaw.ca",
		Url:       "http://www.google.com/thliquor",
		Geo:       Geolocation{Lat: 49.25099, Lon: -122.89683},
		City:      "Burnaby",
		Prov:      "BC",
		Country:   "CA",
		Hour: Hours{
			Open:  [2]int{9, 15},
			Close: [2]int{20, 0}},
		Listing: []Item{
			{Name: "surprise bag",
				Price:       1.99,
				Stock:       3,
				Description: "mix of 2 beers (355ml*2)",
			},
			{Name: "holiday ginks",
				Price:       2.99,
				Stock:       4,
				Description: "a Pinot Noir red wine (750ml imported)",
			},
		},
		Bank: BankAccount{
			Name:    "John Associates",
			Account: "1234567",
			Transit: "5555",
			Inst:    "002",
		},
	}

	insertDoc(Collect{CollectionName: "store", ScopeName: "persona"}, traderJoe, id, cluster)
	// create store's authentification

	type Authenticate_store struct {
		Version    string `json:"_version"`
		Type       string `json:"_type"`
		Created    int64  `json:"_created"` // epoch time in seconds
		CreatedBy  string `json:"_createdBy"`
		Modified   int64  `json:"_modified"` // epoch time in seconds
		ModifiedBy string `json:"_modifiedBy"`
		Id         string `json:"id"`
		User       string `json:"user"`
		Password   string `json:"pass_word"`
		Store_id    string `json:"store_id"`
	}

	counter_id := counter("authenticate_store", cluster)

	hashedPassword, _ := hashcode.HashPassword("timhortons")
	var auth_store_entry Authenticate_store = Authenticate_store{
		Version:   "0.1.0",
		Type:      "authenticate_store",
		Created:   time.Now().Unix(),
		CreatedBy: "class_a_rw",
		Id:        counter_id,
		User:      "joe",
		Password:  hashedPassword,
		Store_id:   "store_" + id,
	}

	insertDoc(Collect{CollectionName: "authenticate_store", ScopeName: "activity"},
		auth_store_entry, counter_id, cluster)

	// create a shopper document

	type PaymentAccount struct {
		Issuer string `json:"issuer"`
		Number string `json:"number"`
		Expiry string `json:"expiry"`
		Cvc    string `json:"cvc"`
	}
	type Shopper struct {
		Version    string         `json:"_version"`
		Type       string         `json:"_type"`
		Created    int64          `json:"_created"` // epoch time in seconds
		CreatedBy  string         `json:"_createdBy"`
		Modified   int64          `json:"_modified"` // epoch time in seconds
		ModifiedBy string         `json:"_modifiedBy"`
		Shopper_id string         `json:"shopper_id"`
		Name       string         `json:"name"`
		Phone      string         `json:"phone"`
		Email      string         `json:"email"`
		Payment    PaymentAccount `json:"payment_account"`
	}
	id = counter("shopper", cluster)
	var shopperAnn = Shopper{
		Version:    "0.1.0",
		Type:       "shopper",
		Created:    time.Now().Unix(),
		CreatedBy:  "class_a_rw",
		Shopper_id: id,
		Name:       "ann",
		Phone:      "2361110000",
		Email:      "ann0000@gmail.com",
		Payment: PaymentAccount{
			Issuer: "visa",
			Number: "4560111122223333",
			Expiry: "08/25",
			Cvc:    "213",
		},
	}

	insertDoc(Collect{CollectionName: "shopper", ScopeName: "persona"}, shopperAnn, id, cluster)

	// create shopper's authentification
	type Authenticate_shopper struct {
		Version    string `json:"_version"`
		Type       string `json:"_type"`
		Created    int64  `json:"_created"` // epoch time in seconds
		CreatedBy  string `json:"_createdBy"`
		Modified   int64  `json:"_modified"` // epoch time in seconds
		ModifiedBy string `json:"_modifiedBy"`
		Id         string `json:"id"`
		User       string `json:"user"`
		Password   string `json:"pass_word"`
		Shopper_id    string `json:"shopper_id"`
	}

	counter_id = counter("authenticate_shopper", cluster)
	hashedPassword, _ = hashcode.HashPassword("annwilson")
	var auth_shopper_entry Authenticate_shopper = Authenticate_shopper{
		Version:   "0.1.0",
		Type:      "authenticate_shopper",
		Created:   time.Now().Unix(),
		CreatedBy: "class_a_rw",
		Id:        counter_id,
		User:      "ann",
		Password:  hashedPassword,
		Shopper_id:   "shopper_" + id,
	}

	insertDoc(Collect{CollectionName: "authenticate_shopper", ScopeName: "activity"},
		auth_shopper_entry, counter_id, cluster)

	// create an admin document
	type Admin struct {
		Version        string `json:"_version"`
		Type           string `json:"_type"`
		Created        int64  `json:"_created"` // epoch time in seconds
		CreatedBy      string `json:"_createdBy"`
		Modified       int64  `json:"_modified"` // epoch time in seconds
		ModifiedBy     string `json:"_modifiedBy"`
		Admin_id       string `json:"admin_id"`
		Name           string `json:"name"`
		Driver_license string `json:"driver_license"`
		Passport       string `json:"passort"`
		City           string `json:"city"`
		Prov           string `json:"prov"`
		Country        string `json:"country"`
		Email          string `json:"email"`
		Phone          string `json:"phone"`
		Class          string `json:"class"`
		Rating         string `json:"rating"`
	}
	id = counter("admin", cluster)
	var admin = Admin{
		Version:        "0.1.0",
		Type:           "admin",
		Created:        time.Now().Unix(),
		CreatedBy:      "class_a_rw",
		Admin_id:       id,
		Name:           "ben",
		Driver_license: "RP77777-1",
		Passport:       "CA-34219-342",
		City:           "Coquitlam",
		Prov:           "BC",
		Country:        "CA",
		Phone:          "2360001111",
		Email:          "ben0000@gmail.com",
		Class:          "b",
		Rating:         "excellent",
	}

	insertDoc(Collect{CollectionName: "admin", ScopeName: "persona"}, admin, id, cluster)

	// create admin's authentification
	type Authenticate_admin struct {
		Version    string `json:"_version"`
		Type       string `json:"_type"`
		Created    int64  `json:"_created"` // epoch time in seconds
		CreatedBy  string `json:"_createdBy"`
		Modified   int64  `json:"_modified"` // epoch time in seconds
		ModifiedBy string `json:"_modifiedBy"`
		Id         string `json:"id"`
		User       string `json:"user"`
		Password   string `json:"pass_word"`
		Admin_id   string `json:"admin_id"`
	}

	counter_id = counter("authenticate_admin", cluster)
	hashedPassword, _ = hashcode.HashPassword("benedward")

	var auth_admin_entry Authenticate_admin = Authenticate_admin{
		Version:   "0.1.0",
		Type:      "authenticate_admin",
		Created:   time.Now().Unix(),
		CreatedBy: "class_a_rw",
		Id:        counter_id,
		User:      "ben",
		Password:  hashedPassword,
		Admin_id:   "admin_" + id,
	}

	insertDoc(Collect{CollectionName: "authenticate_admin", ScopeName: "activity"},
		auth_admin_entry, counter_id, cluster)
}

func counter(collection string, cluster *gocb.Cluster) string {

	counterDoc := collection + "_counter"
	bucket := cluster.Bucket("market")

	collect := bucket.Scope("counter").Collection(counterDoc)
	collect.Binary().Increment(counterDoc, &gocb.IncrementOptions{Initial: 1, Delta: 1})
	result, err := collect.Get(counterDoc, &gocb.GetOptions{})
	if err != nil {
		panic(err)
	}
	var id int
	if err := result.Content(&id); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", id)
}

func insertDoc(coll Collect, v interface{}, id string, cluster *gocb.Cluster) {

	bucket := cluster.Bucket("market")
	collection := bucket.Scope(coll.ScopeName).Collection(coll.CollectionName)
	insertResult, err := collection.Insert(coll.CollectionName+"_"+id, v, nil)
	if err != nil {
		panic(err)
	}
	fmt.Print(coll.CollectionName + " creation success at " + time.Now().String())
	fmt.Println(insertResult)
}

func setup_database(model DataModel, cluster *gocb.Cluster) (*gocb.Cluster, error) {

	// create bucket
	bucketMgr := cluster.Buckets()
	createBucketSettings := gocb.CreateBucketSettings{
		BucketSettings: gocb.BucketSettings{
			Name:       model._Bucket,
			RAMQuotaMB: 200,
			BucketType: gocb.CouchbaseBucketType,
		},
	}
	if err := bucketMgr.CreateBucket(createBucketSettings, &gocb.CreateBucketOptions{}); err != nil {
		return cluster, err
	}

	bucket := cluster.Bucket(model._Bucket)
	collectionMgr := bucket.CollectionsV2()

	// create scopes
	for _, v := range model._Scopes {
		if err := collectionMgr.CreateScope(v, nil); err != nil {
			return cluster, err
		}
	}
	// create collections

	for _, v := range model._Collections {
		if err := collectionMgr.CreateCollection(v.ScopeName, v.CollectionName, nil, nil); err != nil {
			return cluster, err
		}
	}

	return cluster, nil
}