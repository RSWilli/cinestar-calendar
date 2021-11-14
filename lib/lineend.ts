export const crlf = (strings: TemplateStringsArray, ...values: any[]) => {
    return strings.reduce((acc, str, i) => {
        return acc + str.toString() + (i < values.length ? values[i] : '')
    }, "")
        .replace(/\n\s*/gm, "\n")
        .replace(/\n/g, "\n")
}
