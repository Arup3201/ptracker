CREATE VIEW project_summary AS 
SELECT p.id, 
COUNT(t.id) FILTER (WHERE t.status='Unassigned') as unassigned_tasks,
COUNT(t.id) FILTER (WHERE t.status='Ongoing') as ongoing_tasks, 
COUNT(t.id) FILTER (WHERE t.status='Completed') as completed_tasks, 
COUNT(t.id) FILTER (WHERE t.status='Abandoned') as abandoned_tasks
FROM projects as p
LEFT JOIN tasks as t ON p.id=t.project_id
WHERE p.deleted_at IS NULL AND t.deleted_at IS NULL
GROUP BY p.id;