// StickerMaker.java
// Sticker Maker на Java (использует простой алгоритм, без внешних библиотек)

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class StickerMaker {
    public static void main(String[] args) throws IOException {
        String input = null;
        String output = null;
        String colorStr = null;
        int threshold = 30;
        boolean verbose = false;

        for (int i = 0; i < args.length; i++) {
            switch (args[i]) {
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
                    threshold = Integer.parseInt(args[++i]);
                    break;
                case "-v":
                case "--verbose":
                    verbose = true;
                    break;
                case "-h":
                case "--help":
                    System.out.println("Использование: java StickerMaker [опции] <входной_файл>\n" +
                            "  -o, --output <файл>   Выходной файл (PNG)\n" +
                            "  -c, --color <HEX>     Цвет фона\n" +
                            "  -t, --threshold <N>   Порог (0-100)\n" +
                            "  -v, --verbose         Подробный вывод\n" +
                            "  -h, --help            Справка");
                    System.exit(0);
                default:
                    if (input == null) input = args[i];
            }
        }

        if (input == null) {
            System.err.println("Ошибка: не указан входной файл.");
            System.exit(1);
        }

        File inputFile = new File(input);
        if (!inputFile.exists()) {
            System.err.println("Ошибка: файл " + input + " не найден");
            System.exit(1);
        }

        if (output == null) {
            output = input.replaceAll("\\.[^.]+$", "") + "_sticker.png";
        }

        BufferedImage src = ImageIO.read(inputFile);
        int w = src.getWidth();
        int h = src.getHeight();

        // Определяем цвет фона
        Color bgColor = detectBackgroundColor(src);
        if (colorStr != null) {
            try {
                bgColor = Color.decode(colorStr);
            } catch (NumberFormatException e) {
                if (verbose) System.err.println("Предупреждение: неверный формат цвета, используем автоопределение");
            }
        }
        if (verbose) {
            System.out.printf("Цвет фона: RGB(%d,%d,%d)\n", bgColor.getRed(), bgColor.getGreen(), bgColor.getBlue());
        }

        // Создаём изображение с альфа-каналом
        BufferedImage dst = new BufferedImage(w, h, BufferedImage.TYPE_INT_ARGB);
        double maxDist = Math.sqrt(3 * 255 * 255);
        double thresh = threshold / 100.0 * maxDist;

        for (int y = 0; y < h; y++) {
            for (int x = 0; x < w; x++) {
                int rgb = src.getRGB(x, y);
                Color c = new Color(rgb);
                double dist = colorDistance(c, bgColor);
                int alpha = (dist <= thresh) ? 0 : 255;
                dst.setRGB(x, y, (alpha << 24) | (rgb & 0x00FFFFFF));
            }
        }

        ImageIO.write(dst, "png", new File(output));
        if (verbose) {
            System.out.println("Фон удалён. Результат сохранён в " + output);
        }
    }

    private static Color detectBackgroundColor(BufferedImage img) {
        int w = img.getWidth();
        int h = img.getHeight();
        List<Color> pixels = new ArrayList<>();
        // Края
        for (int x = 0; x < w; x++) {
            pixels.add(new Color(img.getRGB(x, 0)));
            pixels.add(new Color(img.getRGB(x, h - 1)));
        }
        for (int y = 0; y < h; y++) {
            pixels.add(new Color(img.getRGB(0, y)));
            pixels.add(new Color(img.getRGB(w - 1, y)));
        }
        // Средний цвет
        long sr = 0, sg = 0, sb = 0;
        for (Color c : pixels) {
            sr += c.getRed();
            sg += c.getGreen();
            sb += c.getBlue();
        }
        int n = pixels.size();
        if (n == 0) return Color.WHITE;
        return new Color((int)(sr / n), (int)(sg / n), (int)(sb / n));
    }

    private static double colorDistance(Color c1, Color c2) {
        double dr = c1.getRed() - c2.getRed();
        double dg = c1.getGreen() - c2.getGreen();
        double db = c1.getBlue() - c2.getBlue();
        return Math.sqrt(dr*dr + dg*dg + db*db);
    }
}
