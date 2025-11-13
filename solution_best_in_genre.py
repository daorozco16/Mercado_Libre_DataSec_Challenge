
#---------- Python 3.12.2 ----------

from typing import Optional
import json
import requests

def bestInGenre(genre: str) -> str:
    """
    Finds the highest-rated TV series in the given genre.

    Parameters:
        genre (str): The genre to search for (e.g., 'Action', 'Comedy', 'Drama')

    Returns:
        str: The name of the highest-rated show in the genre. If there is a tie,
             returns the alphabetically lower name. Returns the name as a string.

    Notes:
        - Ties are broken by alphabetical order of the show name
        - Genre matching is case-insensitive
        - Shows can have multiple genres (comma-separated)
    """
    # Your implementation here

    if not isinstance(genre, str) or not genre.strip():
        # Retorno temprano: si el género no es una cadena de texto o esta vacio o solo tiene espacios, retornamos vacío
        return ""

    base_url = "https://jsonmock.hackerrank.com/api/tvseries"
    page = 1
    total_pages = None

    best_rating: Optional[float] = None
    best_name: Optional[str] = None
    target = genre.strip().lower()

    while True:
        """
        Se podria crear un helper o funcion auxiliar, que haga GET y retorna el JSON como diccionario.
        Y adicional intentar usar requests si está instalado, si no usar urllib.
        """

        resp = requests.get(base_url, params={"page": page}, timeout=10)
        # Se usa para detectar y manejar automaticamente errores HTTP, se puede controlar y excepcionar dentro del Helper
        resp.raise_for_status()
        resp = resp.json()
        items = resp.get("data", []) or []

        for item in items:
            # Extraer géneros y nombre
            item_genres = item.get("genre", "")
            name = item.get("name", "")
            # Si no hay generos o no hay nombre, no se procesa
            if not item_genres or not name:
                continue

            # Se comparan los géneros del string separado por comas y el genero buscado
            genres_list = [g.strip().lower() for g in item_genres.split(",") if g.strip()]
            # Se valida si no pertenece al género buscado
            if target not in genres_list:
                continue

            # Se obtiene el rating y se convierte a float
            rating_val = item.get("imdb_rating", None)
            try:
                rating = float(rating_val)
            except Exception:
                # si no se puede convertir (None, "", "N/A", etc.) lo ignoramos
                continue

            # Actualizar mejor candidato
            if best_rating is None or rating > best_rating:
                best_rating = rating
                best_name = name
            elif rating == best_rating:
                # desempate por orden alfabético (nombre menor)
                if best_name is None or name < best_name:
                    best_name = name

        # ----------------------------------------------------------------------
        # Calcular cantidad total de 'pages'
        if total_pages is None:
            # Validar cuantas paginas existen
            total_pages = resp.get("total_pages", None)
            if total_pages is None:
                # Intentamos obtener total_pages, si no existe se trata de calcular
                # Si vienen 'total' y 'per_page'
                total = resp.get("total", None)
                per_page = resp.get("per_page", None)
                if total is not None and per_page:
                    # Se calcula usando 'per_page - 1' para asegurar que cualquier residuo haga que suba una pagina extra
                    total_pages = (total + per_page - 1) // per_page
                else:
                    # Si definitivamente no existe numero de paginas, implementamos un ultimo mecanismo
                    # Si no existe data en la iteración actual romper el while
                    if not items:
                        break

                    # O forzamos una parada en un límite establecido
                    if page > 100:
                        break

        # Incrementar 'page' o salir cuando se cumpla la cantidad
        if total_pages is not None and page >= total_pages:
            break
        page += 1

    # Se podria implementra una logica mas robusta para poder controlar y capturar los posibles errores, para poderlos retornar:
    #   - No exista genero y no exista nombre en las series
    #   - Si el genero buscado no coincide con ninguna de las series
    #   - Si el formato del rating de la serie seleccionada no se deja convertir

    # Retornar resultado str o cadena vacia si no se encuentra nada
    return best_name if best_name is not None else "No se encontro ninguna serie valida en el genero solicitado!"