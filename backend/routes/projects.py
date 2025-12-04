from flask import Blueprint, request, jsonify
import pydantic

from validation.payload import CreateProjectPayload, CreateTaskPayload, EditTaskPayload, ChangeStatusPayload, ChangeAssigneePayload
from validation.user import User
from services.project import ProjectService
from exceptions import DBOverloadError, DBIntegrityError, NotFoundError, AlreadyExistError
from exceptions.project import NotProjectMemberError, NotProjectOwner, NotTaskAssigneeError

projects_blueprint = Blueprint("projects", __name__)

def list_projects():
    try:
        user = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        projects = ProjectService().list_projects(user_id=user.id)
        return jsonify({
            "message": "fetched all projects", 
            "projects": projects
        })
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Unknown server error occured",
                "details": "We are working on the server, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def create_project():
    try:
        payload = CreateProjectPayload(**request.get_json())
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })

        return jsonify({
            "error": {
                "message": "Input validation failed",
                "details": "Please make sure your input has required fields with their correct type",  
                "errors": errors, 
                "code": "INVALID_INPUT"
            }
        }), 422

    if not payload.name:
        return jsonify({
            "error": {
                "message": "Invalid project name",
                "details": "Field 'name' in project can't be empty",  
                "code": "BAD_REQUEST"
            }
        }), 400
    if not payload.deadline:
        return jsonify({
            "error": {
                "message": "Invalid project deadline",
                "details": "Field 'deadline' in project can't be empty",  
                "code": "BAD_REQUEST"
            }
        }), 400

    try:
        user = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        project = ProjectService().create_projects(name=payload.name, description=payload.description, deadline=payload.deadline, user_id=user.id)
        return jsonify({
            "message": "Project created successfully", 
            "project": project, 
        }), 201
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBIntegrityError as e:
        return jsonify({
            "error": {
                "message": str(e),
                "details": "While creating membership and project intance, integrity error happened",  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def get_project(project_id: str):
    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500
    
    try:
        project = ProjectService().get_project(project_id=project_id, user_id=user_payload.id)
        tasks = ProjectService().get_tasks(project_id=project_id, user_id=user_payload.id)
        return jsonify({
            "project": project, 
            "tasks": tasks
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def get_members(project_id: str):
    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        members, err = ProjectService().get_members(project_id=project_id, user_id=user_payload.id)
        return jsonify({
            "members": members
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def delete_project(project_id: str):
    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        ProjectService().delete_project(project_id=project_id, user_id=user_payload.id)
        return jsonify({
            "message": "Project deleted successfully"
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotProjectOwner as e:
        return jsonify({
            "error": {
                "message": "User is not the project owner",
                "details": str(e),  
                "code": "NOT_OWNER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def join_project(project_code: str):
    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        data = ProjectService().join_project(project_code=project_code, user_id=user_payload.id)
        return jsonify({
            "message": "Successfully joined project", 
            "project": data
        })
    except AlreadyExistError as e:
        return jsonify({
            "error": {
                "message": "User is already a member",
                "details": str(e),  
                "code": "ALREADY_MEMBER"
            }
        }), 409
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def create_task(project_id: str):
    try:
        payload = CreateTaskPayload(**request.get_json())
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })

        return jsonify({
            "error": {
                "message": "Input validation failed",
                "details": "Please make sure your input has required fields with their correct type",  
                "errors": errors, 
                "code": "INVALID_INPUT"
            }
        }), 422

    if not payload.name:
        return jsonify({
            "error": {
                "message": "Invalid task name",
                "details": "Field 'name' in task can't be empty",  
                "code": "BAD_REQUEST"
            }
        }), 400
    if not payload.status:
        return jsonify({
            "error": {
                "message": "Invalid task status",
                "details": "Field 'status' in task can't be empty",  
                "code": "BAD_REQUEST"
            }
        }), 400
    if not payload.assignee:
        return jsonify({
            "error": {
                "message": "Invalid task assignee",
                "details": "Field 'assignee' in task can't be empty",  
                "code": "BAD_REQUEST"
            }
        }), 400

    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        task = ProjectService().create_task(name=payload.name, 
                                     description=payload.description, 
                                     assignee=payload.assignee, 
                                     status=payload.status, 
                                     project_id=project_id, 
                                     user_id=user_payload.id)
        return jsonify({
            "message": "Task created successfully", 
            "task": task 
        }), 201
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotProjectOwner as e:
        return jsonify({
            "error": {
                "message": "User is not a project owner",
                "details": str(e),  
                "code": "NOT_OWNER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def get_task(project_id: str, task_id: str):
    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500
    try:
        task = ProjectService().get_task(task_id=task_id, 
                                         project_id=project_id, 
                                         user_id=user_payload.id)
        return jsonify({
            "task": task 
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def edit_task(project_id: str, task_id: str):
    try:
        payload = EditTaskPayload(**request.get_json())
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })

        return jsonify({
            "error": {
                "message": "Input validation failed",
                "details": "Please make sure your input has required fields with their correct type",  
                "errors": errors, 
                "code": "BAD_REQUEST"
            }
        }), 422

    if not payload.name and not payload.description:
        return jsonify({
            "error": {
                "message": "Invalid edit request",
                "details": "Atleadt one field among 'name' and 'description' need to be present",  
                "code": "BAD_REQUEST"
            }
        }), 400

    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        edited_task = ProjectService().edit_task(task_id=task_id, 
                                   name=payload.name, 
                                   description=payload.description, 
                                   project_id=project_id, 
                                   user_id=user_payload.id)
        return jsonify({
            "message": "Task updated successfully", 
            "task": edited_task
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotTaskAssigneeError as e:
        return jsonify({
            "error": {
                "message": "User is not an owner or the task assignee",
                "details": str(e),  
                "code": "NOT_ASSIGNEE"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def change_task_status(project_id: str, task_id: str):
    try:
        payload = ChangeStatusPayload(**request.get_json())
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })

        return jsonify({
            "error": {
                "message": "Input validation failed",
                "details": "Please make sure your input has required fields with their correct type",  
                "errors": errors, 
                "code": "BAD_REQUEST"
            }
        }), 422

    if not payload.status:
        return jsonify({
            "error": {
                "message": "Invalid parameter to change status request",
                "details": "Field required 'status'",  
                "code": "BAD_REQUEST"
            }
        }), 422

    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        data = ProjectService().change_status(task_id=task_id, 
                                   status=payload.status, 
                                   project_id=project_id, 
                                   user_id=user_payload.id)
        return jsonify({
            "message": "Task status updated successfully", 
            "task": data
        })
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotTaskAssigneeError as e:
        return jsonify({
            "error": {
                "message": "User is not an owner or the task assignee",
                "details": str(e),  
                "code": "NOT_ASSIGNEE"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

def change_assignee(project_id: str, task_id: str):
    try:
        payload = ChangeAssigneePayload(**request.get_json())
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })

        return jsonify({
            "error": {
                "message": "Input validation failed",
                "details": "Please make sure your input has required fields with their correct type",  
                "errors": errors, 
                "code": "BAD_REQUEST"
            }
        }), 422

    if not payload.assignee:
        return jsonify({
            "error": {
                "message": "Invalid parameters to change assignment request",
                "details": "Field required 'assignee'",  
                "code": "BAD_REQUEST"
            }
        }), 422

    try:
        user_payload = User(**request.environ["user"])
    except pydantic.ValidationError as e:
        errors = []
        for err in e.errors():
            errors.append({
                "message": err["msg"], 
                "input": err["input"], 
                "loc": err["loc"]
            })
        print(errors)
        return jsonify({
            "error": {
                "message": "Invalid user data",
                "details": "User data saved at server is corrupted",  
                "code": "SERVER_FAILURE"
            }
        }), 500

    try:
        data = ProjectService().change_assignee(task_id=task_id, 
                                   assignee=payload.assignee, 
                                   project_id=project_id, 
                                   user_id=user_payload.id)
        return jsonify({
            "message": "Task assignee updated successfully", 
            "task": data
        })
    except NotProjectOwner as e:
        return jsonify({
            "error": {
                "message": "User is not an owner of the project",
                "details": str(e),  
                "code": "NOT_OWNER"
            }
        }), 403
    except NotProjectMemberError as e:
        return jsonify({
            "error": {
                "message": "User is not a project member",
                "details": str(e),  
                "code": "NOT_MEMBER"
            }
        }), 403
    except NotFoundError as e:
        return jsonify({
            "error": {
                "message": "Value not found",
                "details": str(e),  
                "code": "NOT_FOUND"
            }
        }), 404
    except DBOverloadError as e:
        return jsonify({
            "error": {
                "message": "Server is overloaded",
                "details": str(e),  
                "code": "SERVER_FAILURE"
            }
        }), 500
    except Exception as e:
        print(str(e))
        return jsonify({
            "error": {
                "message": "Something went wrong in the server",
                "details": "We are working on the error, please try again later",  
                "code": "SERVER_FAILURE"
            }
        }), 500

projects_blueprint.add_url_rule("/", endpoint="list-projects", view_func=list_projects, methods=["GET"])
projects_blueprint.add_url_rule("/", endpoint="create-project", view_func=create_project, methods=["POST"])
projects_blueprint.add_url_rule("/<project_id>", endpoint="get-project", view_func=get_project, methods=["GET"])
projects_blueprint.add_url_rule("/<project_id>", endpoint="delete-project", view_func=delete_project, methods=["DELETE"])
projects_blueprint.add_url_rule("/<project_id>/members", endpoint="get-project-members", view_func=get_members, methods=["GET"])

projects_blueprint.add_url_rule("/join/code/<project_code>", endpoint="join-project", view_func=join_project, methods=["POST"])

projects_blueprint.add_url_rule("/<project_id>/tasks/", endpoint="create-project-task", view_func=create_task, methods=["POST"])
projects_blueprint.add_url_rule("/<project_id>/tasks/<task_id>", endpoint="get-task", view_func=get_task, methods=["GET"])
projects_blueprint.add_url_rule("/<project_id>/tasks/<task_id>", endpoint="edit-project-task", view_func=edit_task, methods=["PUT"])
projects_blueprint.add_url_rule("/<project_id>/tasks/<task_id>/status", endpoint="change-task-status", view_func=change_task_status, methods=["PUT"])
projects_blueprint.add_url_rule("/<project_id>/tasks/<task_id>/assign", endpoint="change-task-assignee", view_func=change_assignee, methods=["PUT"])