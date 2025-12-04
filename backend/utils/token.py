import jwt
from datetime import datetime, timedelta, timezone

from config import Env
from exceptions.auth import JWTError

TOKEN_NAME = "COLLAB_TOKEN"
TOKEN_EXIRES = datetime.utcnow()+timedelta(minutes=30)

def generate_token(payload: dict):
    try:
        token = jwt.encode(payload={
            **payload, 
            "exp": TOKEN_EXIRES, 
            }, key=Env.SECRET_KEY, algorithm="HS256")
        return token
    except jwt.exceptions.InvalidKeyError:
        raise JWTError("Provided key for JWT encoding is invalid")
    except Exception as e:
        print(e)
        raise JWTError("Something went wrong while encoding JWT")

def validate_token(token: str):
    try:
        payload = jwt.decode(jwt=token, key=Env.SECRET_KEY, algorithms="HS256")
        return payload
    except jwt.exceptions.ExpiredSignatureError:
        raise JWTError("JWT token signature has expired")
    except jwt.exceptions.InvalidSignatureError:
        raise JWTError("JWT token signature is invalid") 
    except jwt.exceptions.InvalidTokenError:
        raise JWTError("JWT token is invalid") 
    except jwt.exceptions.InvalidKeyError:
        raise JWTError("JWT token decoding key is invalid") 
    except Exception as e:
        print(e)
        raise JWTError("Something went wrong while decoding JWT") 