<?php
// sticker.php
// Sticker Maker на PHP (использует расширение imagick)

if (php_sapi_name() !== 'cli') {
    die("Это консольное приложение.\n");
}

$options = getopt('o:c:t:vh', ['output:', 'color:', 'threshold:', 'verbose', 'help']);
$args = array_slice($argv, 1);
$input = null;
foreach ($args as $arg) {
    if (!str_starts_with($arg, '-')) {
        $input = $arg;
        break;
    }
}

if (isset($options['h']) || isset($options['help']) || !$input) {
    echo "Использование: php sticker.php [опции] <входной_файл>\n";
    echo "  -o, --output <файл>   Выходной файл (PNG)\n";
    echo "  -c, --color <HEX>     Цвет фона\n";
    echo "  -t, --threshold <N>   Порог (0-100)\n";
    echo "  -v, --verbose         Подробный вывод\n";
    echo "  -h, --help            Справка\n";
    exit(0);
}

if (!file_exists($input)) {
    fwrite(STDERR, "Ошибка: файл $input не найден\n");
    exit(1);
}

$output = $options['o'] ?? $options['output'] ?? null;
if (!$output) {
    $output = preg_replace('/\.[^.]+$/', '', $input) . '_sticker.png';
}

$colorHex = $options['c'] ?? $options['color'] ?? null;
$threshold = isset($options['t']) ? (int)$options['t'] : (isset($options['threshold']) ? (int)$options['threshold'] : 30);
$verbose = isset($options['v']) || isset($options['verbose']);

try {
    $img = new Imagick($input);
    $w = $img->getImageWidth();
    $h = $img->getImageHeight();

    // Определяем цвет фона
    if ($colorHex) {
        $bgColor = new ImagickPixel($colorHex);
    } else {
        // Автоопределение: средний цвет краёв
        $pixels = [];
        // Верхняя и нижняя строки
        for ($x = 0; $x < $w; $x++) {
            $pixels[] = $img->getImagePixelColor($x, 0);
            $pixels[] = $img->getImagePixelColor($x, $h - 1);
        }
        for ($y = 0; $y < $h; $y++) {
            $pixels[] = $img->getImagePixelColor(0, $y);
            $pixels[] = $img->getImagePixelColor($w - 1, $y);
        }
        // Усредняем
        $sr = $sg = $sb = 0;
        foreach ($pixels as $p) {
            $color = $p->getColor();
            $sr += $color['r'];
            $sg += $color['g'];
            $sb += $color['b'];
        }
        $n = count($pixels);
        if ($n > 0) {
            $bgColor = new ImagickPixel(sprintf("rgb(%d,%d,%d)", $sr/$n, $sg/$n, $sb/$n));
        } else {
            $bgColor = new ImagickPixel('white');
        }
    }

    if ($verbose) {
        $c = $bgColor->getColor();
        echo "Цвет фона: RGB({$c['r']},{$c['g']},{$c['b']})\n";
    }

    // Создаём изображение с альфа-каналом
    $img->setImageMatte(true);
    $maxDist = sqrt(3 * 255 * 255);
    $thresh = ($threshold / 100.0) * $maxDist;

    // Проходим по пикселям
    $bgR = $bgColor->getColor()['r'];
    $bgG = $bgColor->getColor()['g'];
    $bgB = $bgColor->getColor()['b'];

    for ($y = 0; $y < $h; $y++) {
        for ($x = 0; $x < $w; $x++) {
            $px = $img->getImagePixelColor($x, $y);
            $c = $px->getColor();
            $dr = $c['r'] - $bgR;
            $dg = $c['g'] - $bgG;
            $db = $c['b'] - $bgB;
            $dist = sqrt($dr*$dr + $dg*$dg + $db*$db);
            if ($dist <= $thresh) {
                $img->setImagePixelColor($x, $y, new ImagickPixel('rgba(0,0,0,0)'));
            } else {
                $img->setImagePixelColor($x, $y, new ImagickPixel("rgba({$c['r']},{$c['g']},{$c['b']},1)"));
            }
        }
    }

    $img->setImageFormat('png');
    $img->writeImage($output);
    $img->destroy();

    if ($verbose) {
        echo "Фон удалён. Результат сохранён в $output\n";
    }
} catch (Exception $e) {
    fwrite(STDERR, "Ошибка: " . $e->getMessage() . "\n");
    exit(1);
}
