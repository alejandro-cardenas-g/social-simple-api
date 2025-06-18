package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/alejandro-cardenas-g/social/internal/store"
)

var usernames = []string{
	"juancarlos",
	"mariapaz",
	"luisfernando",
	"anapaula",
	"carloseduardo",
	"sofiandrea",
	"josemiguel",
	"danielaalex",
	"pedroantonio",
	"angelicamaria",
}

var titles = []string{
	"My First Post",
	"A Day in the Life",
	"Learning Go",
	"Thoughts on Technology",
	"Simple Tips for Productivity",
	"Weekend Vibes",
	"Hello World",
	"Why I Love Coding",
	"Goals for This Month",
	"Quick Update",
	"Working from Home",
	"Favorite Books This Year",
	"Healthy Daily Habits",
	"How I Stay Focused",
	"Morning Routine Ideas",
	"Reflections on the Week",
	"Staying Motivated",
	"Small Wins Matter",
	"What I Learned Today",
	"My Coding Journey",
	"Top 5 Tools I Use",
	"Learning Something New",
	"Mindful Living",
	"Tips for Remote Work",
	"My Workspace Setup",
	"Overcoming Challenges",
	"Notes from the Weekend",
	"Simple Joys of Life",
	"Trying a New Hobby",
	"Looking Ahead",
}

var tags = []string{
	"technology",
	"productivity",
	"lifestyle",
	"personal",
	"motivation",
	"learning",
	"tutorial",
	"tips",
	"travel",
	"health",
	"books",
	"coding",
	"remote-work",
	"daily-life",
	"software",
	"mindset",
	"career",
	"writing",
	"update",
	"habits",
	"focus",
	"creativity",
}

var comments = []string{
	"Great post!",
	"Thanks for sharing!",
	"I totally agree with you.",
	"This was really helpful.",
	"Interesting perspective.",
	"Looking forward to more content like this.",
	"Well explained!",
	"Nice read.",
	"Keep it up!",
	"I learned something new today.",
	"Can you share more details on this?",
	"Love this topic!",
	"Very insightful.",
	"This made my day.",
	"Exactly what I needed.",
	"Clear and concise.",
	"Thanks for the tips!",
	"I will try this out.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)

	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creation user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creation post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creation comment:", err)
			return
		}
	}

	log.Println("Seed executed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)
		users[i] = &store.User{
			Username: username,
			Email:    username + "@mail.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			Content: titles[i%len(titles)] + " Content",
			Title:   titles[i%len(titles)],
			UserId:  user.ID,
			Tags: []string{
				tags[i%len(tags)],
				tags[i%len(tags)],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
