package helpers

var HtmlHeader = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
<style>
        html, body {
            background: #f2f4f8;
            border: 0;
            font-family: Helvetica, Arial, sans-serif;
            font-size: 16px;
            height: 100%;
            margin: 0;
            padding: 0;
        }

        form {
            --background: white;
            --border: rgba(0, 0, 0, 0.125);
            --borderDark: rgba(0, 0, 0, 0.25);
            --borderDarker: rgba(0, 0, 0, 0.5);
            --bgColorH: 0;
            --bgColorS: 0%;
            --bgColorL: 98%;
            --fgColorH: 210;
            --fgColorS: 50%;
            --fgColorL: 38%;
            --shadeDark: 0.3;
            --shadeLight: 0.7;
            --shadeNormal: 0.5;
            --borderRadius: 0.125rem;
            --highlight: #306090;
            background: white;
            border: 1px solid var(--border);
            border-radius: var(--borderRadius);
            box-shadow: 0 1rem 1rem -0.75rem var(--border);
            display: flex;
            flex-direction: column;
            padding: 1rem;
            position: relative;
            overflow: hidden;
        }


        form a:hover {
            color: hsl(var(--fgColorH), var(--fgColorS), calc(var(--fgColorL) * 0.85));
            transition: color 0.25s;
        }

        form a:focus {
            color: hsl(var(--fgColorH), var(--fgColorS), calc(var(--fgColorL) * 0.85));
            outline: 1px dashed hsl(var(--fgColorH), calc(var(--fgColorS) * 2), calc(var(--fgColorL) * 1.15));
            outline-offset: 2px;
        }


        input {
            border: 1px solid var(--border);
            border-radius: var(--borderRadius);
            box-sizing: border-box;
            font-size: 1rem;
            height: 2.25rem;
            line-height: 1.25rem;
            margin-top: 0.25rem;
            order: 2;
            padding: 0.25rem 0.5rem;
            width: 15rem;
            transition: all 0.25s;
        }

        input[type="submit"] {
            color: hsl(var(--bgColorH), var(--bgColorS), var(--bgColorL));
            background: hsl(var(--fgColorH), var(--fgColorS), var(--fgColorL));
            font-size: 0.75rem;
            font-weight: bold;
            margin-top: 0.625rem;
            order: 4;
            outline: 1px dashed transparent;
            outline-offset: 2px;
            padding-left: 0;
            text-transform: uppercase;
        }


    </style>


    <title>pdf_util</title>
</head>
<body>

`

var HtmlHeader2 = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
	<meta http-equiv="content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
<style>
        html, body {
             align-items: center;
            background: #f2f4f8;
            border: 0;
            display: flex;
            font-family: Helvetica, Arial, sans-serif;
            font-size: 16px;
            height: 100%;
            justify-content: center;
            margin: 0;
            padding: 0;
        }

        form {
            --background: white;
            --border: rgba(0, 0, 0, 0.125);
            --borderDark: rgba(0, 0, 0, 0.25);
            --borderDarker: rgba(0, 0, 0, 0.5);
            --bgColorH: 0;
            --bgColorS: 0%;
            --bgColorL: 98%;
            --fgColorH: 210;
            --fgColorS: 50%;
            --fgColorL: 38%;
            --shadeDark: 0.3;
            --shadeLight: 0.7;
            --shadeNormal: 0.5;
            --borderRadius: 0.125rem;
            --highlight: #306090;
            background: white;
            border: 1px solid var(--border);
            border-radius: var(--borderRadius);
            box-shadow: 0 1rem 1rem -0.75rem var(--border);
            display: flex;
            flex-direction: column;
            padding: 1rem;
            position: relative;
            overflow: hidden;
        }


        form a:hover {
            color: hsl(var(--fgColorH), var(--fgColorS), calc(var(--fgColorL) * 0.85));
            transition: color 0.25s;
        }

        form a:focus {
            color: hsl(var(--fgColorH), var(--fgColorS), calc(var(--fgColorL) * 0.85));
            outline: 1px dashed hsl(var(--fgColorH), calc(var(--fgColorS) * 2), calc(var(--fgColorL) * 1.15));
            outline-offset: 2px;
        }


        input {
            border: 1px solid var(--border);
            border-radius: var(--borderRadius);
            box-sizing: border-box;
            font-size: 1rem;
            height: 2.25rem;
            line-height: 1.25rem;
            margin-top: 0.25rem;
            order: 2;
            padding: 0.25rem 0.5rem;
            width: 15rem;
            transition: all 0.25s;
        }

        input[type="submit"] {
            color: hsl(var(--bgColorH), var(--bgColorS), var(--bgColorL));
            background: hsl(var(--fgColorH), var(--fgColorS), var(--fgColorL));
            font-size: 0.75rem;
            font-weight: bold;
            margin-top: 0.625rem;
            order: 4;
            outline: 1px dashed transparent;
            outline-offset: 2px;
            padding-left: 0;
            text-transform: uppercase;
        }


    </style>


    <title>pdf_util</title>
</head>
<body>

`

func main() {

}
