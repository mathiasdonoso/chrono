package report

import "database/sql"

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ReportRepository {
	return &repository{db: db}
}

func (r *repository) GetDailyReport() (Report, error) {
	// Get all task that I been working the previous working day (with the progress)
	// (previous working day is the last day that I worked)

	_ = `SELECT id, name, status, created_at, updated_at
FROM tasks`

	return Report{}, nil

	// s := make([]string, len(statuses))
	// for i, v := range statuses {
	// 	s[i] = fmt.Sprintf("'%s'", v)
	// }
	//
	// query := "SELECT id, name, status, created_at, updated_at FROM tasks WHERE status IN (" + strings.Join(s, ",") + ");"
	//
	// rows, err := r.db.Query(query)
	// if err != nil {
	// 	return []Task{}, err
	// }
	// defer rows.Close()
	//
	// tasks := []Task{}
	// for rows.Next() {
	// 	var task Task
	// 	if err := rows.Scan(
	// 		&task.ID,
	// 		&task.Name,
	// 		&task.Status,
	// 		&task.CreatedAt,
	// 		&task.UpdatedAt,
	// 	); err != nil {
	// 		return []Task{}, err
	// 	}
	// 	tasks = append(tasks, task)
	// }
	//
	// if err := rows.Err(); err != nil {
	// 	return []Task{}, err
	// }
	//
	// return tasks, nil
}
