# Ejercicio Golang Web Server & File Processing.

## Data

Por default, la aplicacion busca los archivos en un directorio llamado "/data".

## Endpoints

**POST /login**

Acepta usuario y contraseña en el cuerpo de la peticion y en formato JSON

Ejemplo:
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"user":"usuario","password":"contraseña"}' \
  http://localhost:3000/login
```  

**GET /files/list**

Se requiere JWT Token en el header y acepta parametro "humanreadable" como query parameter

Ejemplo:

```
curl --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY3NDM3ODAsInVzZXIiOiJ1c3VhcmlvIn0.v919dl32M3WcmYlJcZx2aD3n-OPSOaDWW_GjU7CDXps" \
  --request GET \
  http://localhost:3000/files/list?humanreadable=true
```

**GET /files/metrics**

Se requiere JWT Token en el header y acepta parametro "filename" como query parameter

Ejemplo:

```
curl --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY3NDM3ODAsInVzZXIiOiJ1c3VhcmlvIn0.v919dl32M3WcmYlJcZx2aD3n-OPSOaDWW_GjU7CDXps" \
  --request GET \
  http://localhost:3000/files/metrics?filename=file1.tsv
```
