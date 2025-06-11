# backend_app_tend

Programa Backend para la el proyecto de tendencias.

## Programas Necesarios

### 1: Compilador Golang

<https://go.dev/dl/>

### 2: Base de Datos SQLite

Si posee Windows 10 y 11, y tienes instalado WinGet (<https://apps.microsoft.com/detail/9nblggh4nns1>) puedes tipear en la consola lo siguiente para proceder con la instalación: `winget install --id=SQLite.SQLite  -e`

Otra forma es instalarlo directamente de su sitio oficial: `` luego procederemos a descomprimirlo y dependiendo del sistema, por ejemplo Windows. Procederemos a agregarlo a la ruta "path" de las variables de entorno del sistema.

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

## Pasos para Ejecutar el Backend

### 1. Clonar repositorio de la Rama Dev

Para lograrlo debemos ingresar el comando:
`git clone -b dev https://github.com/Smy250/backend_app_tend`

### 2. Ubicarnos en el directorio raíz del proyecto

Una vez clonado nos ubicamos en el directorio "backend_app_tend". Abrimos el main.go (si se usa visual studio code) o si usas otro ide o editor de texto. Abrimos la terminal y tipearemos el comando `go mod tidy` para que descargue las dependencias del proyecto.

### 3. Ejecución del programa

Finalmente luego de descargar las dependencias, teniendo abierto el archivo main.go en su ide o editor de codigo de preferencia, tipear en la terminal el comando `go run .`

## Extensiones para facilitar el manejo de Golang en Visual Studio Code

### Go - Go Team At Google

### Tooltitude for Go (Golang)
