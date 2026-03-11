export function normalizePhone(raw: string | null | undefined): string {
    if (!raw) return "";
    const trimmed = raw.trim();
    const digits = trimmed.replace(/\D/g, "");
    if (digits.length === 0) return "";

    let cc = "30";
    let national = digits;

    if (trimmed.startsWith('+')) {
        if (digits.length > 10) {
            cc = digits.slice(0, digits.length - 10);
            national = digits.slice(-10);
        } else {
            // no explicit country code digits provided, fallback to gr
            cc = "30";
            national = digits.slice(-10);
        }
    } else {
        if (digits.length > 10) {
            cc = digits.slice(0, digits.length - 10);
            national = digits.slice(-10);
        } else {
            cc = "30";
            national = digits;
        }
    }

    const part1 = national.slice(0, 3);
    const part2 = national.slice(3);
    return `+${cc} ${part1}${part2 ? ' ' + part2 : ''}`.trim();
}

export function isValidPhone(raw: string | null | undefined): boolean {
    if (!raw) return false;
    const trimmed = raw.trim();
    const digits = trimmed.replace(/\D/g, "");
    if (digits.length < 7) return false;

    // If digits length > 10, assume leading digits are country code
    const national = digits.length > 10 ? digits.slice(-10) : digits;

    return national.length === 10;
}
