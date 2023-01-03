package database

import (
	"context"

	"github.com/go-redis/redis/v9"
)

/*
Uses Redis as data store. Go redis client- go-redis.
*/

/*
NoSQL Schema

Redis Data Layout:

	Hashes-

		1) key:- user:[id]

			(user:[id] field value ...)

			type User struct stored as field value pairs keyed by "user:[User.id]"
			id = INCR next_user_id (see 4)


		2) key:- users

			(users name id)

			Secondary index to 1)


		3) key:- auths

			(auths auth-token id)

			Maps auth token to user id for authentication.
			Also auth is stored as a field of 1.


	Strings-

		4) key:- next_user_id

			(next_user_id value)

			Incremented every time new user created. value used as User.ID


	Sorted Sets

		5) key:- followers:[id]

			(followers:[id] UNIXtime id ...)

			Each key (read as "followers of [user:id]") represents a user whose followers are contained in its sorted set
			sorted by the UNIX time at which the user was followed by the user represented by the id value.


		6) key:- following:[id]

			(following:[id] UNIXtime id ...)

			Similar to 5, but the reverse. Sorted set of user:[id] and the users followed by them.

		note: when a user 1 follows a user 2, add user 2 to 6) and user 1 to 5).


	GeoSpatial Indices

		7) key:- locations

			(locations longitude latitude [id])

			GeoHash or sorted set of geo-coordinate locations of users.
			Used to find nearby friends.

			note-	Consider partitioning key space into countries/regions to
				  	improve performance as size of indices M increases.
				 	Efficiency of O(N + log M) for search by radius.

	TODO-   Add set/sorted set user:auths:id to keep track of multiple logins and auth tokens.
*/

// shortcut to pass ctx generally.
var ctx context.Context = context.TODO()

// func NewClient creates and returns a new instance of a redis client connected with
// the redis server specified by opt.
// Also checks if the database server is actually responding to requests.
//
// TODO: Leverage connections pool effectively & Transport Layer Security
func NewClient(opt *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
