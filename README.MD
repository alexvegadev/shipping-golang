
# Desafío técnico ML

  

## Ejecución del proyecto

### Requerimientos:

**Postgres >= 10**

  

**Golang >= 1.16.6**

### Ejecución

**Si se va a hacer una ejecución local, es necesario crear la db "ml_shipping" en postgres y ejecutar el archivo schema.sql para crear la tabla.**

Clonar el repositorio: **git clone** https://github.com/alexvegadev/shipping-golang.git

Crear un archivo **.env** se puede setear con las siguientes configuraciones:
- POSTGRES_URL=postgresql://localhost/ml_shipping?user=postgres&password=root&sslmode=disable
- PORT=5000
- DEPOSIT_LIMIT=100

A través de la consola, usar el comando: **go run main.go**
  

## Requerimiento

### Objetivo

El Negocio de MercadoLibre se sigue expandiendo y ahora decidió incursionar en el almacenamiento de los productos de los vendedores. Para ese fin, nos fue requerido construir un microservicio para mantener un control de stock de los productos. Esta aplicación debe mantener la relación entre los productos publicados en el marketplace y su ubicación dentro de cada depósito.

### Requerimientos

1. Exponer un endpoint REST para agregar productos en una ubicación.

	2. Se nos indicará el Depósito, producto, cantidad y ubicación donde quiere colocar.

	3. Validar que la dirección tenga el patrón correcto.

	4. Que el producto/item sea almacenado en nuestros depósitos.

	5. No se pueden colocar más de 3 productos distintos en una ubicación.

	6. La suma de las cantidades de los productos que hubiera en una ubicación no puede ser mayor a 100 unidades.

2. Exponer un endpoint para poder retirar productos de una ubicación. Se nos indicará el depósito, producto, cantidad y ubicación de donde sacarla.

3. Exponer un endpoint de lectura. Se nos indica un depósito y una ubicación, y este liste los productos y cantidad que hay en el mismo.

4. Exponer un endpoint de búsqueda. Se nos indica el depósito y producto, y este nos devuelva las posibles ubicaciones y cantidad en las mismas.

  

## Detalle de la solución

Para hacer las consultas a la API de mercado libre se utilizaron las librerías http, ioutil y json también se crearon DTO's en donde solo se traen los campos que son imprescindibles, se crearon modelos para manejar el depósito y se incluyeron las validaciones especificadas para la dirección del depósito.
  

## Deploy
En este caso no se hizo deploy, pero se puede probar localmente.

## Documentación

Ya que el modelo está basado en el proyecto realizado en java, la estructura no cambia.

Se puede usar la misma documentación: https://mldesafio.herokuapp.com/swagger-ui.html
