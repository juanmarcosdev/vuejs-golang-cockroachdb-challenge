# Challenge (Reto)
![Challenge official image](https://i.ibb.co/rMjNVqn/Captura-de-pantalla-de-2020-06-25-11-31-45.png)
## Application Architecture (Arquitectura de la Aplicación)
![Application Architecture image](https://i.ibb.co/9pRv0Bc/Captura-de-pantalla-de-2020-06-25-11-43-31.png)
### Español
Esta aplicación está diseñada con una arquitectura de **microservicios**, construida gracias al uso de contenedores de Docker, quiere decir que para poder correrla sólo tienes que tener instalada en tu máquina Docker y Docker Compose (Windows, macOS o Linux).
#### Descarga e instalación
Para poder correr la aplicación debes:

 1. Tener instalado Docker y Docker Compose en tu máquina (de ahora en adelante, máquina host). Aquí puedes encontrar información de cómo instalar:
	 - Windows: [Instalar Docker en Windows](https://docs.docker.com/docker-for-windows/install/)
	 - macOS: [Instalar Docker en macOS](https://docs.docker.com/docker-for-mac/install/)
	 - Linux: [Instalar Docker en Linux](https://docs.docker.com/engine/install/), [Instalar Docker en Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/como-instalar-y-usar-docker-en-ubuntu-18-04-1-es)
	 - Docker Compose: [Instalar Docker Compose en Windows, macOS o Linux](https://docs.docker.com/compose/install/) (dentro de la web, seleccionar la pestaña del sistema operativo correspondiente)
2. Clonar o descargar este repositorio:
	* `git clone https://github.com/juanmarcos935/vuejs-golang-cockroachdb-challenge.git` (desde la consola)
	*  O puedes descargarlo accediendo a la pestaña Clone > elegir la opción "Download ZIP". Debes extraerlo y obtener la carpeta de nombre  `vuejs-golang-cockroachdb-challenge`.
3. Tras tener la carpeta lista, entramos en ella:
	* `cd vuejs-golang-cockroachdb-challenge`
4. Y finalmente debemos ejecutar el comando que levantará todos los 3 contenedores correspondientes a cada microservicio (Frontend, Backend & Database):
	* `docker-compose up --build`
	* El flag `--build` es útil en cuanto a que construirá de cero las imágenes que poseen Dockerfile anexo.
	* Nota: Para que todo funcione correctamente, debes asegurarte de que en la máquna host (tu máquina) estén libres los puertos:
		* 3000, donde se accede a la página web (Frontend)
		* 5000, donde se encuentra el servidor corriendo (Backend)
		* 26257, donde corre el contenedor de persistencia (Base de Datos)
5. Cuando terminen de levantarse todos los contenedores, ¡ya está listo! puedes acceder a la aplicación desde la página web poniendo localhost:3000 en tu web browser (te sugerimos que uses Google Chrome). 
	* También asegúrate de contar con una conexión a internet, debido a que la aplicación realiza peticiones  una API externa y esto requiere de conexión  a internet.
### English
This application is designed with the software architecture of **microservices**, built thanks to the use of Docker containers, it means that you can run it having installed only Docker and Docker Compose (Windows, macOS o Linux) on your machine.
#### Download an installation
In order to run the application you need to:

 1. Have Docker and Docker Compose in your local machine (host machine). Here you can find information about it:
	 - Windows: [Install Docker on Windows](https://docs.docker.com/docker-for-windows/install/)
	 - macOS: [Install Docker on macOS](https://docs.docker.com/docker-for-mac/install/)
	 - Linux: [Install Docker on Linux](https://docs.docker.com/engine/install/), [Install Docker on Ubuntu 18.04](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-18-04)
	 - Docker Compose: [Install Docker Compose on Windows, macOS or Linux](https://docs.docker.com/compose/install/) (inside the page, choose the tab that corresponds with your operative system)
2. Clone or download this repository:
	* `git clone https://github.com/juanmarcos935/vuejs-golang-cockroachdb-challenge.git` (from the console)
	*  Or you can download it clicking on the tab Clone > choose the option "Download ZIP". You need to extract it and get the folder called  `vuejs-golang-cockroachdb-challenge`.
3. Having the folder ready, we move inside it:
	* `cd vuejs-golang-cockroachdb-challenge`
4. And finally to execute the command that will deploy the 3 containers from every microservice (Frontend, Backend & Database):
	* `docker-compose up --build`
	* The flag `--build` is useful due to it will build the images from scratch (the microservices that depend on a specific Dockerfile)
	* Annotation: In order to everything work as expected, you need to be sure that the following ports are not used by other processes:
		* 3000, where you can access to the web page (Frontend)
		* 5000, where the server is running (Backend)
		* 26257, where the persistance container is running (Database)
5. When all three containers are up, ¡everything is ready! you can access the application's web page typing localhost:3000 in your web browser (we encourage you to use Google Chrome). 
	* Be sure you have an stable internet connection, because the application sends requests to an extern API and this requieres internet connectivity.
