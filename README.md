# Challenge (Reto)
![Challenge official image](https://i.ibb.co/rMjNVqn/Captura-de-pantalla-de-2020-06-25-11-31-45.png)
## Application Architecture (Arquitectura de la Aplicación)
![Application Architecture image](https://i.ibb.co/9pRv0Bc/Captura-de-pantalla-de-2020-06-25-11-43-31.png)
* [Ir a la Documentación en Español](#1-spanish)
* [Go to English Documentation](#2-english)
<h2> ¿Qué aprendí llevando a cabo este proyecto? </h2>
<ul>
	<li>Conceptos de Vue.js: componentes, templates, SPA (single page application), distintas directivas de Vue como v-bind, v-router, v-if, v-for, etc</li>
	<li>Uso del framework Tailwind CSS</li>
	<li>Construcción de backend (un servidor) con el lenguaje de Go, usando fasthttp y como router fasthttprouter</li>
	<li>Conceptos de Go como: Handlers en Servidor, Structs, Receiver Functions, Funciones públicas y privadas, uso de librería para JSON encoding/json, uso de librería para ejecutar comandos del sistema operativo os/exec, uso de la librería strings, time, sort</li>
	<li>Web Scrapping con Go mediante librerías goquery y x/net/html</li>
	<li>Conexión de Servidor Go con Base de Datos SQL mediante librerías database/sql y lib/pq</li>
	<li>Creación y manejo de cluster simple y multi-nodo de el SGBD CockroachDB</li>
	<li>Manejo de una aplicación sencilla mediante microservicios usando despliegue y orquestación con Docker y Docker Compose</li>
</ul>
<h2> What did I learn doing this project? </h2>
<ul>
	<li>Vue.js concepts: components, templates, SPA (single page application), different Vue directives such as v-bind, v-router, v-if, v-for, etc.</li>
	<li>Use of the Tailwind CSS framework.</li>
	<li>Construction of backend (a server) with the language of Go, using fasthttp and as router fasthttprouter.</li>
	<li>Go concepts such as: Handlers in Server, Structs, Receiver Functions, Public and private functions, use of library encoding/json for JSON usage, use of library os/exec to execute operating system commands, use of library strings, time, sort</li>
	<li>Web Scrapping with Go using goquery and x/net/ html libraries</li>
	<li>Go Server connection with SQL Database through database/sql and lib/pq libraries</li>
	<li>Creation and management of simple and multi-node cluster of the DBMS CockroachDB</li>
	<li>Simple application management through microservices using deployment and orchestration with Docker and Docker Compose.</li>
</ul>
<h3 id="1-spanish">Documentación en Español</h3>
Esta aplicación está diseñada con una arquitectura de **microservicios**, construida gracias al uso de contenedores de Docker, quiere decir que para poder correrla sólo tienes que tener instalada en tu máquina Docker y Docker Compose (Windows, macOS o Linux)
<h4>Descarga e instalación</h4>
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
		* 8080, donde se accede a la página web (Frontend)
		* 5000, donde se encuentra el servidor corriendo (Backend)
		* 26257, donde corre el contenedor de persistencia (Base de Datos)
5. Cuando terminen de levantarse todos los contenedores, ¡ya está listo! puedes acceder a la aplicación desde la página web poniendo localhost:3000 en tu web browser (te sugerimos que uses Google Chrome). 
	* También asegúrate de contar con una conexión a internet, debido a que la aplicación realiza peticiones  una API externa y esto requiere de conexión  a internet.
	* Nota: Existe un cuarto (4) contenedor llamado db-builder. Este contenedor es temporal y no hace parte de la arquitectura final de la aplicación, pues su única función es manipular el contenedor de la base de datos (database) e inicializar el esquema de la base de datos. Al finalizar esta tarea, el contenedor muere y solo sobreviven los 3 principales.
#### Uso
* Nota: Para el frontend se utilizó el framework Vue.js junto con el preprocesador de CSS, Tailwind CSS
##### Homepage
Ya habiéndose desplegado todos los 3 contenedores, puede ingresar a su navegador a la dirección http://localhost:8080 y verá una página similar a esta:

![homepage](https://i.ibb.co/pzwxHt1/Captura-de-pantalla-2020-06-28-03-32-45.png)

Este es el Home, o la página principal. Posee una pequeña descripción acerca del reto. Tienes dos enlaces abajo: uno para dirigirte al Endpoint #1 y otro para acceder al Endpoint #2. 

##### Endpoint #1

Haciendo click hacia el Endpoint #1, verás esta pantalla: 

![endpoint1](https://i.ibb.co/HDNdrvB/Captura-de-pantalla-2020-06-28-03-43-04.png)
Aquí tienes:
* Un botón para Realizar consulta.
* Un enlace para volver a la página de Inicio (Home).

Si haces click sobre el botón de Realizar consulta vas a ver lo siguiente:

![query](https://i.ibb.co/YWqJrXJ/Captura-de-pantalla-2020-06-28-03-46-21.png)

En el campo de texto, puedes colocar el dominio a consultar en el formato 'web.com', como por ejemplo, google.com o github.com.
Si no especificas un dominio la página lanzará error:
![error](https://i.ibb.co/2KvDdBb/Captura-de-pantalla-2020-06-28-03-47-51.png)

Sin embargo, si se especifica uno válido, y se le da click al botón Enviar, se hará la solicitud y mientras tanto habrá un loader:

![loader](https://i.ibb.co/cQCFZRV/Captura-de-pantalla-2020-06-28-03-50-11.png)
Tras haber tenido éxito la petición, se muestra el mensaje:
![exito](https://i.ibb.co/S7CxvgG/Captura-de-pantalla-2020-06-28-03-51-43.png)
Dando click en OK, se puede ver la información de la solicitud:
![registro1](https://i.ibb.co/Zz64y3j/Captura-de-pantalla-2020-06-28-03-53-17.png)
* Nota 1: Si aparece un valor de "null" en el campo de Anterior Grado SSL, esto se debe a que no existía un registro previo de hace una hora o menos (50 min, 40 min, 30 min, 1 min, etc.) que tuviera información del dominio especificado. Por ejemplo, si se hace una solicitud casi inmediatamente después esta vez tendrá la información correspondiente al primer registro y se mostrará el valor correcto:
![registro2](https://i.ibb.co/mJrTSYn/Captura-de-pantalla-2020-06-28-03-55-10.png)
* Nota 2: El web scrapping correspondiente a la búsqueda del enlace o link del logo se hace de tal forma que busca en el HTML del dominio en su etiqueta html -> luego la etiqueta head -> después la etiqueta link -> y busca la etiqueta que tenga como atributo "rel" con el valor de "shortcut icon", y devuelve el valor que posea el atributo "href" de la misma etiqueta. Como en ocasiones los sitios web están estructurados de manera distinta y guardan sus logos en una estructura que no es ésta, no es posible obtener el link del logo, por ello se muestra un campo vacío en la columna correspondiente.
* Nota 3: También puede suceder que sea un dominio que no sea muy consultado en la API de SSL Labs, y esto conlleve a que cuando se realiza la petición a la API externa de dicha empresa apenas estén realizando los cálculos y tests de seguridad para determinar el grado de sus servidores. También puede ocurrir que el enlace del logo de una página no apunte a un enlace web, sino a una ruta relativa del sistema de archivos del servidor local del dominio. Ambas sucesos suceden en la siguiente consulta:
![uao](https://i.ibb.co/pyBJBXQ/Captura-de-pantalla-2020-06-28-03-59-40.png)
* Nota 4: También puede darse el caso de que sea imposible obtener el Owner de los servidores:
![fb](https://i.ibb.co/GPP1j8Z/Captura-de-pantalla-2020-06-28-04-05-24.png)
Sin embargo, la herramienta es capaz de obtener información sobre dominios que tienen muchos servidores:
![netflix](https://i.ibb.co/1frm7Kc/Captura-de-pantalla-2020-06-28-04-07-54.png)
Y también casos en los que se logra obtener toda la información:
![youtube](https://i.ibb.co/tsM5yTJ/Captura-de-pantalla-2020-06-28-04-20-26.png)
##### Endpoint #2
Haciendo click hacia el Endpoint #2, verás esta pantalla: 
![endpoint2](https://i.ibb.co/643zqsW/Captura-de-pantalla-2020-06-28-04-11-43.png)
Donde sencillamente se tiene a disposición el botón de Consultar dominios, que mostrará los dominios que han sido consultados en el Endpoint #1. Esto maneja persistencia con la base de datos, así que es posible cerrar el navegador por completo y volver a la página y seguirán estando los dominios.
Al dar click en el botón se muestra:
![endp2](https://i.ibb.co/txfFtN6/Captura-de-pantalla-2020-06-28-04-13-37.png)
 Si damos click en el botón Desplegar y aún no habíamos hecho consultas en el Endpoint #1, nos arroja error indicándolo:
![nodominios](https://i.ibb.co/5jrHhmB/Captura-de-pantalla-2020-06-28-03-48-49.png)
Sin embargo, si se ha consultado 1 o más, notifica que hubo éxito:
![exito2](https://i.ibb.co/BsqLqPK/Captura-de-pantalla-2020-06-28-04-17-07.png)
Y lista a todos los dominios:
![listadominios](https://i.ibb.co/xqRGc5t/Captura-de-pantalla-2020-06-28-04-18-18.png)
<h3 id="2-english">English Documentation</h3>
This application is designed with the software architecture of **microservices**, built thanks to the use of Docker containers, it means that you can run it having installed only Docker and Docker Compose (Windows, macOS o Linux) on your machine
<h4>Download and installation</h4>
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
		* 8080, where you can access to the web page (Frontend)
		* 5000, where the server is running (Backend)
		* 26257, where the persistance container is running (Database)
5. When all three containers are up, ¡everything is ready! you can access the application's web page typing localhost:3000 in your web browser (we encourage you to use Google Chrome). 
	* Be sure you have an stable internet connection, because the application sends requests to an extern API and this requieres internet connectivity.
	* Annotation: There exists a fourth (4th) container called db-builder. This is a temporary container and doesn't make part of the final application architecture, his main functionality is to manipulate the database container and initialize the database through bash and SQL scripts. When it finishes this task, the container dies and only survives the other 3.
#### Usage
* Annotation: The framework Vue.js was used to build the frontend, alongside the CSS preprocessor Tailwind CSS
##### Homepage
With the 3 containers all up, you can access to the address http://localhost:8080 in your web browser and you'll see a page like this:

![homepage](https://i.ibb.co/pzwxHt1/Captura-de-pantalla-2020-06-28-03-32-45.png)

This is the homepage. It has a little description about the challenge. There are two links: one for accessing to the Endpoint #1 and the other for the access to Endpoint #2. 

##### Endpoint #1

If you go to Endpoint #1, you'll see this: 

![endpoint1](https://i.ibb.co/HDNdrvB/Captura-de-pantalla-2020-06-28-03-43-04.png)
You got:
* A button to do the query
* A link to go back to home

If you click the query button you will see:

![query](https://i.ibb.co/YWqJrXJ/Captura-de-pantalla-2020-06-28-03-46-21.png)

In the textfield, you can write the domain that is going to search, using the format 'web.com', for example, google.com or github.com.
If you don't specify a valid domain it will throw an error:
![error](https://i.ibb.co/2KvDdBb/Captura-de-pantalla-2020-06-28-03-47-51.png)

However, if you specific a valid one, and you click on the Send button, the request will be made and there will be a loader for the time it gets the response:

![loader](https://i.ibb.co/cQCFZRV/Captura-de-pantalla-2020-06-28-03-50-11.png)
If the request is successful, it shows the message:
![exito](https://i.ibb.co/S7CxvgG/Captura-de-pantalla-2020-06-28-03-51-43.png)
If you click OK, you can see the information below:
![registro1](https://i.ibb.co/Zz64y3j/Captura-de-pantalla-2020-06-28-03-53-17.png)
* Annotation 1: If there is a "null" value in the Previous SSL Grade, this is because there was no register in the database that is from a hour ago or less (50 min, 40 min, 30 min, 1 min, etc.) that had the information of the specified domain. For example, if you do a query immediately after the first query it will show the right information:
![registro2](https://i.ibb.co/mJrTSYn/Captura-de-pantalla-2020-06-28-03-55-10.png)
* Annotation 2: The web scrapping that belongs to the search of the logo link it functions on a certain way: it searches in the HTML of the webpage for an html node -> then for head node -> finally link node -> and searches for the node that has "rel" property with value "shortcut icon", and returns the value of "href" of the self node that fits the condition. Because all webpages have different structure of saving their icon or logo, if a site doesn't have this structure it won't show any value in logo field. 
* Annotation 3: It can happen too that if a domain is not very popular and because of that not very often searched in SSL Labs API, results that when the request is due SSL Labs server is still running security tests on the domain to determine the SSL grade of each server, so it can came with blank values on the grades of the servers. It can happen too that the link to the logo of the webpage is not a webpage link, but a relative path of the filesystem that belongs to the server running the webpage. Here we got a case where happens these two things:
![uao](https://i.ibb.co/pyBJBXQ/Captura-de-pantalla-2020-06-28-03-59-40.png)
* Annotation 4: Sometimes it will be imposible to retrieve the data from the Owners of the servers:
![fb](https://i.ibb.co/GPP1j8Z/Captura-de-pantalla-2020-06-28-04-05-24.png)
However, the web tool is capable of getting data from a domain that has a lot of servers:
![netflix](https://i.ibb.co/1frm7Kc/Captura-de-pantalla-2020-06-28-04-07-54.png)
And of course it is capable of retrieving all data required without leaving any blank value:
![youtube](https://i.ibb.co/tsM5yTJ/Captura-de-pantalla-2020-06-28-04-20-26.png)
##### Endpoint #2
If you click to Endpoint #2 you'll see:
![endpoint2](https://i.ibb.co/643zqsW/Captura-de-pantalla-2020-06-28-04-11-43.png)
Where we got only a button that you can query all the Domains you have searched for in Endpoint #1. This is handled with persistency using a database, so you can close your browser, return to the page and you still got the information of the queried domains.
If you click the button shows up this:
![endp2](https://i.ibb.co/txfFtN6/Captura-de-pantalla-2020-06-28-04-13-37.png)
 If we click on Show and yet there are no domains that we have searched for in Endpoint #1, it will throw an error saying that:
![nodominios](https://i.ibb.co/5jrHhmB/Captura-de-pantalla-2020-06-28-03-48-49.png)
However if we have already searched one or more domains, it says that was successful:
![exito2](https://i.ibb.co/BsqLqPK/Captura-de-pantalla-2020-06-28-04-17-07.png)
And lists all the domains:
![listadominios](https://i.ibb.co/xqRGc5t/Captura-de-pantalla-2020-06-28-04-18-18.png)
