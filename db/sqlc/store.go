package frog_blossom_db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store Functions for executing db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a db transaction
func (store *Store) executeTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// InitSetupConfigTx populates db tables with initial site-specific config data
// Use user info to populate the `posts`, `pages`, and `meta` tables
func (store *Store) InitSetupConfigTx(ctx context.Context, args InitSetupConfigTxParams) (InitSetupConfigTxResult, error) {
	var result InitSetupConfigTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUsers(ctx, args.UserId)
		if err != nil {
			return fmt.Errorf("get users err: %v", err)
		}
		result.User = user

		for _, postsParams := range args.InitialPosts {
			post, err := q.CreatePosts(ctx, postsParams)
			if err != nil {
				return fmt.Errorf("create posts err: %v", err)
			}
			result.Posts = append(result.Posts, post)
		}

		for _, pageParams := range args.InitialPages {
			page, err := q.CreatePages(ctx, pageParams)
			if err != nil {
				return fmt.Errorf("create pages err: %v", err)
			}
			result.Pages = append(result.Pages, page)
		}

		for _, metaParas := range args.InitialMeta {
			meta, err := q.CreateMeta(ctx, metaParas)
			if err != nil {
				return fmt.Errorf("create meta err: %v", err)
			}
			result.Metas = append(result.Metas, meta)
		}

		return nil
	})
	return result, err
}

// CreatePostsTx creates new posts content based on user information
// It utilizes user info(users.id, users.username) to create the `posts` and its respective `meta`.
func (store *Store) CreatePostsTx(ctx context.Context, args CreatePostTxParams) (CreatePostTxResult, error) {
	var result CreatePostTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUsers(ctx, args.UserId)
		if err != nil {
			return fmt.Errorf("get users err: %v", err)
		}

		postArgs := CreatePostsParams{
			Title:        args.Posts.Title,
			Content:      args.Posts.Content,
			AuthorID:     user.ID,
			Url:          args.Posts.Url,
			Status:       args.Posts.Status,
			PublishedAt:  time.Time{},
			EditedAt:     time.Time{},
			PostAuthor:   user.Username,
			PostMimeType: args.Posts.PostMimeType,
			PublishedBy:  user.Username,
			UpdatedBy:    user.Username,
		}

		result.Posts, err = q.CreatePosts(ctx, postArgs)
		if err != nil {
			return fmt.Errorf("create posts err: %v", err)
		}

		metaArgs := CreateMetaParams{
			PostsID:         sql.NullInt64{Int64: result.Posts.ID, Valid: true},
			MetaTitle:       args.Metas.MetaTitle,
			MetaDescription: args.Metas.MetaDescription,
			MetaRobots:      args.Metas.MetaRobots,
			MetaOgImage:     args.Metas.MetaOgImage,
			Locale:          args.Metas.Locale,
			PageAmount:      args.Metas.PageAmount,
			SiteLanguage:    args.Metas.SiteLanguage,
			MetaKey:         args.Metas.MetaKey,
			MetaValue:       args.Metas.MetaValue,
		}

		result.Metas, err = q.CreateMeta(ctx, metaArgs)
		if err != nil {
			return fmt.Errorf("create meta err: %v", err)
		}

		return nil
	})
	return result, err
}

// CreatePageTx creates new pages content based on user information
// It utilizes user info(users.id, users.username) to create the `page` and its respective `meta`.
func (store *Store) CreatePageTx(ctx context.Context, args CreatePageTxParams) (CreatePageTxResult, error) {
	var result CreatePageTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUsers(ctx, args.UserId)
		if err != nil {
			return fmt.Errorf("get users err: %v", err)
		}

		pageArgs := CreatePagesParams{
			Domain:         args.Pages.Domain,
			AuthorID:       user.ID,
			PageAuthor:     user.Username,
			Title:          args.Pages.Title,
			Url:            args.Pages.Url,
			MenuOrder:      args.Pages.MenuOrder,
			ComponentType:  args.Pages.ComponentType,
			ComponentValue: args.Pages.ComponentValue,
			PageIdentifier: args.Pages.PageIdentifier,
			OptionID:       args.Pages.OptionID,
			OptionName:     args.Pages.OptionName,
			OptionValue:    args.Pages.OptionValue,
			OptionRequired: args.Pages.OptionRequired,
		}

		result.Pages, err = q.CreatePages(ctx, pageArgs)
		if err != nil {
			return fmt.Errorf("create pages err: %v", err)
		}

		metaArgs := CreateMetaParams{
			PageID:          sql.NullInt64{Int64: result.Pages.ID},
			MetaTitle:       args.Metas.MetaTitle,
			MetaDescription: args.Metas.MetaDescription,
			MetaRobots:      args.Metas.MetaRobots,
			MetaOgImage:     args.Metas.MetaOgImage,
			Locale:          args.Metas.Locale,
			PageAmount:      args.Metas.PageAmount,
			SiteLanguage:    args.Metas.SiteLanguage,
			MetaKey:         args.Metas.MetaKey,
			MetaValue:       args.Metas.MetaValue,
		}

		result.Metas, err = q.CreateMeta(ctx, metaArgs)
		if err != nil {
			return fmt.Errorf("create meta err: %v", err)
		}

		return nil
	})

	return result, err
}

// UpdatePostsTx updates existing content in the `posts` table and its respective `meta` table.
// It utilizes user info (users.id, users.username) to update the content and its associated metadata.
func (store *Store) UpdatePostsTx(ctx context.Context, args UpdatePostTxParams) (UpdatePostTxResult, error) {
	var result UpdatePostTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUsers(ctx, args.UserId)
		if err != nil {
			return fmt.Errorf("get user err: %v", err)
		}

		post, err := q.GetPosts(ctx, *args.PostId)
		if err != nil {
			return fmt.Errorf("get post err: %v", err)
		}

		meta, err := q.GetMetaByPostsIDForUpdate(ctx, sql.NullInt64{Int64: *args.PostId, Valid: true})
		if err != nil {
			return fmt.Errorf("get meta err: %v", err)
		}

		postArgs := UpdatePostsParams{
			ID:           post.ID,
			Title:        args.Posts.Title,
			Content:      args.Posts.Content,
			AuthorID:     user.ID,
			Url:          args.Posts.Url,
			UpdatedAt:    args.Posts.UpdatedAt,
			Status:       args.Posts.Status,
			PublishedAt:  args.Posts.PublishedAt,
			EditedAt:     args.Posts.EditedAt,
			PostAuthor:   user.Username,
			PostMimeType: args.Posts.PostMimeType,
			PublishedBy:  args.Posts.PublishedBy,
			UpdatedBy:    args.Posts.UpdatedBy,
		}

		metaArgs := UpdateMetaParams{
			ID:              meta.ID,
			PostsID:         sql.NullInt64{Int64: post.ID, Valid: true},
			MetaTitle:       args.Metas.MetaTitle,
			MetaDescription: args.Metas.MetaDescription,
			MetaRobots:      args.Metas.MetaRobots,
			MetaOgImage:     args.Metas.MetaOgImage,
			Locale:          args.Metas.Locale,
			PageAmount:      args.Metas.PageAmount,
			SiteLanguage:    args.Metas.SiteLanguage,
			MetaKey:         args.Metas.MetaKey,
			MetaValue:       args.Metas.MetaValue,
		}

		result.Posts, err = q.UpdatePosts(ctx, postArgs)
		if err != nil {
			return fmt.Errorf("update post err: %v", err)
		}

		result.Metas, err = q.UpdateMeta(ctx, metaArgs)
		if err != nil {
			return fmt.Errorf("update meta err: %v", err)
		}

		return nil
	})
	return result, err
}

// UpdatePageTx updates existing content in the `page` table and its respective `meta` table.
// It utilizes user info (users.id, users.username) to update the content and its associated metadata.
func (store *Store) UpdatePageTx(ctx context.Context, args UpdatePageTxParams) (UpdatePageTxResult, error) {
	var result UpdatePageTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUsers(ctx, args.UserId)
		if err != nil {
			return fmt.Errorf("get user err: %v", err)
		}

		page, err := q.GetPages(ctx, *args.PageId)
		if err != nil {
			return fmt.Errorf("get pages err: %v", err)
		}

		meta, err := q.GetMetaByPageIDForUpdate(ctx, sql.NullInt64{Int64: *args.PageId, Valid: true})
		if err != nil {
			return fmt.Errorf("get meta err: %v", err)
		}

		pageArgs := UpdatePagesParams{
			ID:             page.ID,
			Domain:         args.Pages.Domain,
			AuthorID:       user.ID,
			PageAuthor:     user.Username,
			Title:          args.Pages.Title,
			Url:            args.Pages.Url,
			MenuOrder:      args.Pages.MenuOrder,
			ComponentType:  args.Pages.ComponentType,
			ComponentValue: args.Pages.ComponentValue,
			PageIdentifier: args.Pages.PageIdentifier,
			OptionID:       args.Pages.OptionID,
			OptionName:     args.Pages.OptionName,
			OptionValue:    args.Pages.OptionValue,
			OptionRequired: args.Pages.OptionRequired,
		}

		metaArgs := UpdateMetaParams{
			ID:              meta.ID,
			PageID:          sql.NullInt64{Int64: page.ID, Valid: true},
			MetaTitle:       args.Metas.MetaTitle,
			MetaDescription: args.Metas.MetaDescription,
			MetaRobots:      args.Metas.MetaRobots,
			MetaOgImage:     args.Metas.MetaOgImage,
			Locale:          args.Metas.Locale,
			PageAmount:      args.Metas.PageAmount,
			SiteLanguage:    args.Metas.SiteLanguage,
			MetaKey:         args.Metas.MetaKey,
			MetaValue:       args.Metas.MetaValue,
		}

		result.Pages, err = q.UpdatePages(ctx, pageArgs)
		if err != nil {
			return fmt.Errorf("update pages err: %v", err)
		}

		result.Metas, err = q.UpdateMeta(ctx, metaArgs)
		if err != nil {
			return fmt.Errorf("update meta err: %v", err)
		}

		return nil
	})
	return result, err
}

// DeletePostsTx deletes existing content in the `posts` table and its respective `meta` table.
// It utilizes posts info (post.id) to delete posts content and its associated metadata.
func (store *Store) DeletePostsTx(ctx context.Context, args DeletePostTxParams) (DeletePostTxResult, error) {
	var result DeletePostTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		err = q.DeleteMetaByPostId(ctx, sql.NullInt64{
			Int64: *args.PostId,
			Valid: true,
		})
		if err != nil {
			return fmt.Errorf("delete meta err: %v", err)
		}
		result.DeletedMeta = true

		err = q.DeletePosts(ctx, *args.PostId)
		if err != nil {
			return fmt.Errorf("delete posts err: %v", err)
		}
		result.DeletedPost = true

		return nil
	})
	return result, err
}

// DeletePageTx deletes existing content in the `page` table and its respective `meta` table.
// It utilizes posts info (post.id) to delete posts content and its associated metadata.
func (store *Store) DeletePageTx(ctx context.Context, args DeletePageTxParams) (DeletePageTxResult, error) {
	var result DeletePageTxResult

	err := store.executeTx(ctx, func(q *Queries) error {
		var err error

		err = q.DeleteMetaByPageId(ctx, sql.NullInt64{
			Int64: *args.PageId,
			Valid: true,
		})
		if err != nil {
			return fmt.Errorf("delete meta err: %v", err)
		}
		result.DeletedMeta = true

		err = q.DeletePages(ctx, *args.PageId)
		if err != nil {
			return fmt.Errorf("delete page err: %v", err)
		}
		result.DeletedPage = true

		return nil
	})
	return result, err
}
