package mongo

import (
	"github.com/AIGPTku/api-aigptku.id/app/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionInvitationDetail = "invitation_detail"
)

type repoMongo struct {
	master *mongo.Database
	trx *mongo.Database
}

func New(dbMaster, dbTrx *mongo.Database) repository.MongoInterface {
	return &repoMongo{
		master: dbMaster,
		trx: dbTrx,
	}
}