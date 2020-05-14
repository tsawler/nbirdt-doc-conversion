package main

import (
	"context"
	"fmt"
	"github.com/gosimple/slug"
	"strings"
	"time"
)

// HoldingFile describes holding file model
type HoldingFile struct {
	ID                int `json:"id"`
	HoldingNameEn     string
	HoldingID         int       `json:"holding_id"`
	DisplayNameEn     string    `json:"display_name_en"`
	DisplayNameFr     string    `json:"display_name_fr"`
	FileDescriptionEn string    `json:"file_description_en"`
	FileDescriptionFr string    `json:"file_description_fr"`
	FileNameDisplay   string    `json:"file_name_display"`
	FileName          string    `json:"file_name"`
	Active            int       `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// ProjectFile describes project file model
type ProjectFile struct {
	ID                int `json:"id"`
	ProjectNameEn     string
	ProjectID         int       `json:"project_id"`
	DisplayNameEn     string    `json:"display_name_en"`
	DisplayNameFr     string    `json:"display_name_fr"`
	FileDescriptionEn string    `json:"file_description_en"`
	FileDescriptionFr string    `json:"file_description_fr"`
	FileNameDisplay   string    `json:"file_name_display"`
	FileName          string    `json:"file_name"`
	Active            int       `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// PublicationFile describes publication file model
type PublicationFile struct {
	ID                int `json:"id"`
	PublicationNameEn string
	PublicationID     int       `json:"publication_id"`
	DisplayNameEn     string    `json:"display_name_en"`
	DisplayNameFr     string    `json:"display_name_fr"`
	FileDescriptionEn string    `json:"file_description_en"`
	FileDescriptionFr string    `json:"file_description_fr"`
	FileNameDisplay   string    `json:"file_name_display"`
	FileName          string    `json:"file_name"`
	Active            int       `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (app *application) addSlugToHoldings() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, holding_name_en from holdings order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		s := slug.Make(name)
		slugValue := fmt.Sprintf("%s-%d", s, id)
		stmt := "update holdings set slug = $1 where id = $2"
		_, err = app.db.ExecContext(ctx, stmt, slugValue, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (app *application) getAllHoldingDocs() ([]HoldingFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select hf.id, hf.holding_id, hf.display_name_en, hf.display_name_fr,
			hf.file_description_en, hf.file_description_fr, hf.file_name, hf.active,
			hf.created_at, hf.updated_at, hf.file_name_display, h.holding_name_en
			
			from holding_files hf
			left join holdings h on (h.id = hf.holding_id)
			order by holding_id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var files []HoldingFile

	for rows.Next() {
		var i HoldingFile
		err = rows.Scan(
			&i.ID,
			&i.HoldingID,
			&i.DisplayNameEn,
			&i.DisplayNameFr,
			&i.FileDescriptionEn,
			&i.FileDescriptionFr,
			&i.FileName,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FileNameDisplay,
			&i.HoldingNameEn,
		)
		if err != nil {
			return files, err
		}
		files = append(files, i)
	}

	return files, nil
}

func (app *application) updateFileNamesForHoldings() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, file_name_display from holding_files order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		stmt := "update holding_files set file_name = $1 where id = $2"

		oldDisplayName := name
		last4 := oldDisplayName[len(oldDisplayName)-4:]
		rootName := strings.TrimSuffix(oldDisplayName, last4)
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)

		_, err = app.db.ExecContext(ctx, stmt, newFileName, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (app *application) addSlugToPublications() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, publication_name_en from publications order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		s := slug.Make(name)
		slugValue := fmt.Sprintf("%s-%d", s, id)
		stmt := "update publications set slug = $1 where id = $2"
		_, err = app.db.ExecContext(ctx, stmt, slugValue, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (app *application) getAllPublicationDocs() ([]PublicationFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select hf.id, hf.publication_id, hf.display_name_en, hf.display_name_fr,
			hf.file_description_en, hf.file_description_fr, hf.file_name, hf.active,
			hf.created_at, hf.updated_at, hf.file_name_display, h.publication_name_en
			
			from publication_files hf
			left join publications h on (h.id = hf.publication_id)
			order by publication_id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var files []PublicationFile

	for rows.Next() {
		var i PublicationFile
		err = rows.Scan(
			&i.ID,
			&i.PublicationID,
			&i.DisplayNameEn,
			&i.DisplayNameFr,
			&i.FileDescriptionEn,
			&i.FileDescriptionFr,
			&i.FileName,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FileNameDisplay,
			&i.PublicationNameEn,
		)
		if err != nil {
			return files, err
		}
		files = append(files, i)
	}

	return files, nil
}

func (app *application) updateFileNamesForPublications() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, file_name_display from publication_files order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}

		stmt := "update publication_files set file_name = $1 where id = $2"

		oldDisplayName := name
		last4 := oldDisplayName[len(oldDisplayName)-4:]
		rootName := strings.TrimSuffix(oldDisplayName, last4)
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)

		_, err = app.db.ExecContext(ctx, stmt, newFileName, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (app *application) addSlugToProjects() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, project_name_en from projects order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		s := slug.Make(name)
		slugValue := fmt.Sprintf("%s-%d", s, id)
		stmt := "update projects set slug = $1 where id = $2"
		_, err = app.db.ExecContext(ctx, stmt, slugValue, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (app *application) getAllProjectDocs() ([]ProjectFile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select hf.id, hf.project_id, hf.display_name_en, hf.display_name_fr,
			hf.file_description_en, hf.file_description_fr, hf.file_name, hf.active,
			hf.created_at, hf.updated_at, hf.file_name_display, h.project_name_en
			
			from project_files hf
			left join projects h on (h.id = hf.project_id)
			order by project_id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var files []ProjectFile

	for rows.Next() {
		var i ProjectFile
		err = rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.DisplayNameEn,
			&i.DisplayNameFr,
			&i.FileDescriptionEn,
			&i.FileDescriptionFr,
			&i.FileName,
			&i.Active,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FileNameDisplay,
			&i.ProjectNameEn,
		)
		if err != nil {
			return files, err
		}
		files = append(files, i)
	}

	return files, nil
}

func (app *application) updateFileNamesForProjects() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// get docs
	q := `select id, file_name_display from project_files order by id`

	rows, err := app.db.QueryContext(ctx, q)

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		stmt := "update project_files set file_name = $1 where id = $2"

		oldDisplayName := name
		last4 := oldDisplayName[len(oldDisplayName)-4:]
		rootName := strings.TrimSuffix(oldDisplayName, last4)
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)

		_, err = app.db.ExecContext(ctx, stmt, newFileName, id)
		if err != nil {
			return err
		}
	}

	return nil

}
