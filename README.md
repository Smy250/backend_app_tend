# backend_app_tend

Programa Backend para el proyecto de tendencias.

## Programas Necesarios

### 1: Compilador Golang

<https://go.dev/dl/>

Como dato adicional es caso de obtener: "`failed to initialize database, got error Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub`". Será necesario instalar MSYS2 especificamente los compiladores de C/C++. Primeramente accediendo al siguiente enlace: "<https://www.msys2.org/>".

Una vez en el sitio entraremos iremos bajando en el sitio hasta ver el apartado:

`Installation`

`Download the installer: msys2-x86_64-20250221.exe`

Procederemos a descargarlo. Haciendo click en el nombre del archivo, luego lo instalaremos.

De ser necesario, visualizar las variables de entorno posterior a la instalación (en caso de tener Windows "Editar las variables de entorno").

Si no se encuentra la siguiente carpeta en path: `C:\msys64\mingw64\bin` agregarla.

Una vez comprobado todo en orden en el menu de inicio buscar el programa "`MSYS2 MSYS`". Una vez abierto veremos la terminal y en él colocaremos los siguientes comandos:

* `pacman -Su`
* `"pacman -S --needed base-devel mingw-w64-x86_64-toolchain" o "pacman -S base-devel mingw-w64-x86_64-toolchain"`

Pulsaremos Enter o Y, dependiendo el caso. Esperamos que finalice y finalmente estaremos listos para seguir con los pasos.

En caso de que se nos queje por la falta de SQLite, se procedería a instalarlo directamente en la terminal MSYS2 con el siguiente comando: `pacman -S mingw-w64-x86_64-sqlite3` o también podrías instalarlo desde su sitio oficial: `https://www.sqlite.org/download.html` luego procederemos a descomprimirlo y dependiendo del sistema, por ejemplo Windows. Procederemos a agregarlo a la ruta "path" de las variables de entorno del sistema, si usaste pacman para instalarlo no es necesario especificar la ruta ya que estará en la misma ruta del compilador de C++.

## Pasos para Ejecutar el Backend

### 1. Clonar repositorio de la Rama principal

Para lograrlo debemos ingresar el comando:
`git clone https://github.com/Smy250/backend_app_tend`

### 2. Ubicarnos en el directorio raíz del proyecto

Una vez clonado nos ubicamos en el directorio "backend_app_tend". Abrimos el main.go (si se usa visual studio code) o si usas otro ide o editor de texto. Abrimos la terminal y tipearemos el comando `go mod tidy` para que descargue las dependencias del proyecto.

### 3. Ejecución del programa

Finalmente luego de descargar las dependencias, teniendo abierto el archivo main.go en su ide o editor de codigo de preferencia, tipear en la terminal el comando `go run .` o si prefieres tener el ejecutable `go build .`

## Extensiones para facilitar el manejo de Golang en Visual Studio Code

### Go - Go Team At Google

### Tooltitude for Go (Golang)
