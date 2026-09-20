export const formatRupiah = (value: number) =>
  "Rp " + new Intl.NumberFormat("id-ID").format(value);

/** Compact form for axis labels/summaries, e.g. "Rp 86,4 jt" for 86_400_000. */
export const formatRupiahCompact = (value: number) => {
  if (Math.abs(value) >= 1_000_000) {
    const millions = value / 1_000_000;
    const rounded = Math.round(millions * 10) / 10;
    return "Rp " + rounded.toLocaleString("id-ID") + " jt";
  }
  return formatRupiah(value);
};

export const formatShortDate = (iso: string) =>
  new Date(iso).toLocaleDateString("en-GB", { day: "numeric", month: "short" });

export const formatWeekRange = (weekStartIso: string) => {
  const start = new Date(weekStartIso);
  const end = new Date(start);
  end.setDate(start.getDate() + 6);

  const startMonth = start.toLocaleDateString("en-GB", { month: "short" });
  const endMonth = end.toLocaleDateString("en-GB", { month: "short" });

  if (startMonth === endMonth) {
    return `${start.getDate()}–${end.getDate()} ${endMonth}`;
  }
  return `${start.getDate()} ${startMonth}–${end.getDate()} ${endMonth}`;
};
