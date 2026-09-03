// StickerMaker.cs
// Sticker Maker на C# (использует SixLabors.ImageSharp)

using SixLabors.ImageSharp;
using SixLabors.ImageSharp.PixelFormats;
using SixLabors.ImageSharp.Processing;
using System;
using System.Collections.Generic;
using System.IO;

class StickerMaker
{
    static void Main(string[] args)
    {
        string input = null, output = null, colorStr = null;
        int threshold = 30;
        bool verbose = false;

        for (int i = 0; i < args.Length; i++)
        {
            switch (args[i])
            {
                case "-o":
                case "--output":
                    output = args[++i];
                    break;
                case "-c":
                case "--color":
                    colorStr = args[++i];
                    break;
                case "-t":
                case "--threshold":
                    threshold = int.Parse(args[++i]);
                    break;
                case "-v":
                case "--verbose":
                    verbose = true;
                    break;
                case "-h":
                case "--help":
                    Console.WriteLine(@"Использование: dotnet run -- [опции] <входной_файл>
  -o, --output <файл>   Выходной файл (PNG)
  -c, --color <HEX>     Цвет фона
  -t, --threshold <N>   Порог (0-100)
  -v, --verbose         Подробный вывод
  -h, --help            Справка");
                    return;
                default:
                    if (input == null) input = args[i];
                    break;
            }
        }

        if (input == null)
        {
            Console.Error.WriteLine("Ошибка: не указан входной файл.");
            Environment.Exit(1);
        }
        if (!File.Exists(input))
        {
            Console.Error.WriteLine($"Ошибка: файл {input} не найден");
            Environment.Exit(1);
        }

        output = output ?? input.Replace(Path.GetExtension(input), "") + "_sticker.png";

        using (var image = Image.Load<Rgba32>(input))
        {
            int w = image.Width, h = image.Height;

            // Определяем цвет фона
            Rgba32 bgColor = DetectBackgroundColor(image);
            if (colorStr != null)
            {
                try
                {
                    bgColor = ColorHexToRgba(colorStr);
                }
                catch
                {
                    if (verbose) Console.Error.WriteLine("Предупреждение: неверный формат цвета, используем автоопределение");
                }
            }
            if (verbose)
            {
                Console.WriteLine($"Цвет фона: RGB({bgColor.R},{bgColor.G},{bgColor.B})");
            }

            double maxDist = Math.Sqrt(3 * 255 * 255);
            double thresh = threshold / 100.0 * maxDist;

            // Создаём новое изображение с альфа-каналом
            using (var dst = new Image<Rgba32>(w, h))
            {
                for (int y = 0; y < h; y++)
                {
                    for (int x = 0; x < w; x++)
                    {
                        Rgba32 px = image[x, y];
                        double dist = ColorDistance(px, bgColor);
                        byte alpha = (dist <= thresh) ? (byte)0 : (byte)255;
                        dst[x, y] = new Rgba32(px.R, px.G, px.B, alpha);
                    }
                }
                dst.SaveAsPng(output);
            }
        }

        if (verbose)
        {
            Console.WriteLine($"Фон удалён. Результат сохранён в {output}");
        }
    }

    private static Rgba32 DetectBackgroundColor(Image<Rgba32> img)
    {
        int w = img.Width, h = img.Height;
        var pixels = new List<Rgba32>();
        for (int x = 0; x < w; x++)
        {
            pixels.Add(img[x, 0]);
            pixels.Add(img[x, h - 1]);
        }
        for (int y = 0; y < h; y++)
        {
            pixels.Add(img[0, y]);
            pixels.Add(img[w - 1, y]);
        }
        long sr = 0, sg = 0, sb = 0;
        foreach (var p in pixels)
        {
            sr += p.R;
            sg += p.G;
            sb += p.B;
        }
        int n = pixels.Count;
        if (n == 0) return new Rgba32(255, 255, 255);
        return new Rgba32((byte)(sr / n), (byte)(sg / n), (byte)(sb / n));
    }

    private static double ColorDistance(Rgba32 c1, Rgba32 c2)
    {
        double dr = c1.R - c2.R;
        double dg = c1.G - c2.G;
        double db = c1.B - c2.B;
        return Math.Sqrt(dr * dr + dg * dg + db * db);
    }

    private static Rgba32 ColorHexToRgba(string hex)
    {
        hex = hex.TrimStart('#');
        if (hex.Length == 6)
        {
            byte r = Convert.ToByte(hex.Substring(0, 2), 16);
            byte g = Convert.ToByte(hex.Substring(2, 2), 16);
            byte b = Convert.ToByte(hex.Substring(4, 2), 16);
            return new Rgba32(r, g, b);
        }
        throw new FormatException("Неверный формат HEX цвета");
    }
}
