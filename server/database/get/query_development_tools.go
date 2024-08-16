package get

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"server/models"
)

func DevelopmentToolsQuery(database *pgx.Conn) []models.DevelopmentTools {
	// Query command
	query, err := database.Query(context.Background(), "SELECT tools.id, tools.field, tools.descriptions, tools.text_color, tools.span, tools.order, icons.path, icons.icon_names FROM favourite_tools as tools INNER JOIN favourite_tools_icons as icons on tools.field = icons.field ORDER BY tools.id")
	if err != nil {
		log.Err(err).Msg("[Database] Error querying development tools: %v\n")
	}

	// Iterate through the get results and assign it into returned value
	var returnedLiteral []models.DevelopmentTools
	for query.Next() {
		// Holding get values for each field
		var field, description, textColor, path string
		var iconNames []string
		var id, span, order int

		err := query.Scan(&id, &field, &description, &textColor, &span, &order, &path, &iconNames)
		if err != nil {
			log.Err(err).Msg("[Database] Scan error: %v\n")
			return nil
		}

		returnedLiteral = append(returnedLiteral, models.DevelopmentTools{
			Id:          id,
			Field:       field,
			Description: description,
			Style:       models.Style{Span: span, TextColor: textColor},
			Icons:       models.Icons{Path: path, Name: iconNames},
			Order:       order,
		})
	}

	// Close get result on function return
	defer query.Close()

	return returnedLiteral
}
