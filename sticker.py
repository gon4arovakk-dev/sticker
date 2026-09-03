# sticker.py
# Sticker Maker на Python (использует rembg)

import sys
import os
import argparse
from PIL import Image
from rembg import remove

def main():
    parser = argparse.ArgumentParser(description="Sticker Maker (Background Removal) - Python")
    parser.add_argument("input", help="Входной файл")
    parser.add_argument("-o", "--output", help="Выходной файл (PNG)")
    parser.add_argument("-c", "--color", help="Цвет фона (HEX) - не используется в rembg, оставлен для совместимости")
    parser.add_argument("-t", "--threshold", type=int, default=30, help="Порог (не используется)")
    parser.add_argument("-v", "--verbose", action="store_true", help="Подробный вывод")
    args = parser.parse_args()

    if not os.path.exists(args.input):
        print(f"Ошибка: файл {args.input} не найден")
        sys.exit(1)

    output = args.output or os.path.splitext(args.input)[0] + "_sticker.png"

    try:
        with open(args.input, 'rb') as i:
            input_data = i.read()
            output_data = remove(input_data)
            with open(output, 'wb') as o:
                o.write(output_data)
        if args.verbose:
            print(f"Фон удалён. Результат сохранён в {output}")
    except Exception as e:
        print(f"Ошибка обработки: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
