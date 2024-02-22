# Proyecto BasicGORM
Este proyecto es una librería para construir y ejecutar consultas SQL de forma dinámica utilizando Golang.

## Descripción
BasicGORM es un paquete en Go diseñado para facilitar la interacción con bases de datos SQL utilizando el paquete database/sql. Además de ejecutar consultas SQL de forma programática utilizando métodos encadenados para establecer diferentes partes de la consulta, como la cláusula SELECT, WHERE, JOIN, ORDER BY, etc., BasicGORM proporciona funciones para validar los datos antes de realizar operaciones de inserción, actualización y eliminación en la base de datos


## Características

- Ejecución de Consultas SQL: Permite ejecutar consultas SQL de forma fácil y segura en bases de datos SQL compatibles..
- Soporte para diferentes cláusulas SQL como SELECT, WHERE, JOIN, ORDER BY, etc.
- Posibilidad de ejecutar consultas personalizadas.
- Manejo de errores integrado.
- Validación de Datos: Antes de realizar operaciones de inserción, actualización o eliminación en la base de datos, BasicGORM valida los datos proporcionados según los criterios definidos en el esquema de la tabla. Esto incluye verificaciones de longitud mínima o máxima, tipo de datos correcto, valores permitidos, entre otros.
- Transacciones: Admite transacciones SQL para agrupar operaciones en una única unidad de trabajo, lo que garantiza la consistencia y la integridad de los datos en escenarios complejos que requieren varias operaciones.

## Instalación

Para instalar BasicGORM, simplemente ejecuta:

```bash
go get github.com/deybin/basicgorm

```
## Uso

 - variables de entorno

	 - ENV_DDBB_SERVER='IP/Host'
	 - ENV_DDBB_USER='user'
	 - ENV_DDBB_PASSWORD='password'
	 - ENV_DDBB_DATABASE='database_default'
	 - ENV_DDBB_PORT=port
	 - ENV_DDBB_SSL ='false/true'
	 - ENV_KEY_CRYPTO='key_para_encriptar'

```bash
package main

import (
	"fmt"
	"github.com/deybin/basicgorm"
)

func main() {
	query := basicgorm.Querys{
		Table: "mi_tabla",
	}

	result, err := query.Select().Where("campo", basicgorm.I, "valor_campo").And("campo2", basicgorm.IN, []interface{}{"valor_filtro_IN1", "valor_filtro_IN2", "valor_filtro_IN3"}).Exec(basicgorm.QConfig{Database: "mi_database"}).OrderBy("campo_ordenar DESC").All()
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(result)
}
```

## Contribución
¡Las contribuciones son bienvenidas! Si quieres contribuir a este proyecto o encuentras algún problema por favor abre un issue primero para discutir los cambios propuestos.

## Licencia
© Deybin, 2024~time.Now

Este proyecto está bajo la Licencia MIT. Consulta el archivo [LICENSE](https://github.com/deybin/basicgorm/master/LICENSE) para más detalles.