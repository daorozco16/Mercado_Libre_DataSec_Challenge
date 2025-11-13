
#---------- Python 3.12.2 ----------

def validar_tablero(board: list) -> bool:
    """
    Validar que board tenga una estructura correcta.

    - Debe ser una lista 2D es decir lista de listas
    - No puede estar vacio
    - Todas las filas deben tener el mismo número de columnas
    - Cada celda solo puede contener 0 o 1
    """

    # Verificar que sea una lista y no este vacia
    if not isinstance(board, list) or not board:
        raise ValueError("El tablero debe ser una lista y no estar vacia.")

    # Verificar que cada fila sea una lista
    for fila in board:
        if not isinstance(fila, list):
            raise ValueError("Cada fila del tablero debe ser una lista.")

    # Verificar que todas las filas tengan el mismo número de columnas
    columnas = len(board[0])
    for i, fila in enumerate(board):
        if len(fila) != columnas:
            raise ValueError(f"La fila {i} tiene {len(fila)} columnas, pero se esperaban {columnas}.")


    # Verificar que los valores sean solo 0 o 1
    for i, fila in enumerate(board):
        for j, celda in enumerate(fila):
            if celda not in (0, 1):
                raise ValueError(f"La celda ({i},{j}) contiene un valor invalido: {celda}. Solo se permiten 0 o 1.")

def count_neighbouring_mines(board: list) -> list:
    """
    Counts neighbouring mines for each cell in a Minesweeper board.

    Parameters:
        board (list): A 2D list where 0 represents an empty space and 1 represents a mine

    Returns:
        list: A 2D list where each cell contains the count of neighbouring mines,
              or 9 if the cell contains a mine
    """
    # Your implementation here

    #Validar board antes de procesarlo
    validar_tablero(board)

    # Comprender el tamaño de la matriz
    rows = len(board)
    cols = len(board[0])

    # Crear una matriz del mismo tamaño usando comprensión de listas
    result = [[0 for _ in range(cols)] for _ in range(rows)]


    """
    Contar las celdas vecinas usando coordenadas relativas
        Posiciones cercanas a la celda actual
        Medidas en relación a su ubicación

    Cada celda tiene hasta 8 posibles vecinos
        Arriba, abajo, izquierda, derecha y 4 diagonales
    """
    directions = [
        (-1, -1), (-1, 0), (-1, 1),  # fila anterior
        (0, -1),           (0, 1),   # misma fila
        (1, -1),  (1, 0),  (1, 1)    # fila siguiente
    ]

    # Analizar cada una de las celdas de la matriz
    for i in range(rows):
        for j in range(cols):

            # Si la celda es una mina se pone 9
            if board[i][j] == 1:
                result[i][j] = 9

            # Si la celda no es una mina, contar las minas vecinas usando las coordenadas relativas
            else:
                count = 0 # Contador de las posiciones cercanas que son minas
                for dx, dy in directions:
                    ni, nj = i + dx, j + dy  # Nueva posición que se usara para asegurar limites

                    # Se asegura que las coordenadas existen y estan dentro de los limites
                    if 0 <= ni < rows and 0 <= nj < cols:

                        # Si la celda es una mina, se suma 1 al contador
                        if board[ni][nj] == 1:
                            count += 1

                result[i][j] = count

    return result