"""Generate the static figures embedded in README.md from data/Housing.csv.

Run from the repository root:  python docs/scripts/make_figures.py
Requires: pandas, matplotlib.  The running application does NOT use Python.
"""

from pathlib import Path

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "docs" / "images"
DATA = ROOT / "data" / "Housing.csv"

ACCENT = "#22d3ee"
ACCENT2 = "#a78bfa"
PINK = "#f472b6"
LIME = "#a3e635"
MUTED = "#8ea0bb"
GRID = "#243149"

plt.rcParams.update(
    {
        "figure.facecolor": "#0b1020",
        "axes.facecolor": "#111a2b",
        "axes.edgecolor": "#2a3852",
        "axes.labelcolor": "#c7d3e6",
        "axes.titlecolor": "#e7eefb",
        "axes.titlesize": 12,
        "axes.titleweight": "bold",
        "xtick.color": MUTED,
        "ytick.color": MUTED,
        "text.color": "#e7eefb",
        "grid.color": GRID,
        "grid.linestyle": "--",
        "grid.alpha": 0.6,
        "font.size": 10,
    }
)


def head(fig, title, subtitle):
    fig.suptitle(title, x=0.055, y=0.94, ha="left", fontsize=19, fontweight="bold")
    fig.text(0.056, 0.895, subtitle, ha="left", color=MUTED, fontsize=10)


def save(fig, name):
    OUT.mkdir(parents=True, exist_ok=True)
    fig.savefig(OUT / name, dpi=150, bbox_inches="tight", facecolor=fig.get_facecolor())
    plt.close(fig)
    print("wrote", OUT / name)


def distributions(df):
    numeric = df.select_dtypes(include="number").columns.tolist()
    fig, axes = plt.subplots(2, 3, figsize=(14, 7.5))
    for ax, column, color in zip(axes.flat, numeric, [ACCENT, ACCENT2, LIME, PINK, "#fbbf24", "#60a5fa"]):
        values = df[column].dropna()
        ax.hist(values, bins=18, color=color, alpha=0.82, edgecolor="#0b1020", linewidth=0.6)
        ax.axvline(values.mean(), color="#ffffff", linestyle="--", linewidth=1.2, alpha=0.75)
        ax.set_title(column)
        ax.set_xlabel(f"mean {values.mean():,.0f}   σ {values.std():,.0f}", color=MUTED, fontsize=9)
        ax.set_ylabel("count")
        ax.grid(axis="y")
    fig.subplots_adjust(left=0.06, right=0.97, bottom=0.09, top=0.83, hspace=0.42, wspace=0.25)
    head(fig, "01 — Numeric distributions", "Histogram + mean marker for every numeric feature · Housing.csv")
    save(fig, "01-distributions.png")


def correlation(df):
    numeric = df.select_dtypes(include="number")
    corr = numeric.corr()
    fig, ax = plt.subplots(figsize=(9, 7.4))
    im = ax.imshow(corr, cmap="coolwarm", vmin=-1, vmax=1)
    ax.set_xticks(range(len(corr)), corr.columns, rotation=42, ha="right")
    ax.set_yticks(range(len(corr)), corr.columns)
    for i in range(len(corr)):
        for j in range(len(corr)):
            ax.text(j, i, f"{corr.iloc[i, j]:.2f}", ha="center", va="center", fontsize=9,
                    color="#ffffff" if abs(corr.iloc[i, j]) > 0.55 else MUTED)
    ax.set_xticks(np.arange(-0.5, len(corr), 1), minor=True)
    ax.set_yticks(np.arange(-0.5, len(corr), 1), minor=True)
    ax.grid(which="minor", color="#0b1020", linewidth=2)
    ax.tick_params(which="minor", bottom=False, left=False)
    for spine in ax.spines.values():
        spine.set_color(GRID)
    cbar = fig.colorbar(im, ax=ax, fraction=0.046, pad=0.04)
    cbar.ax.yaxis.set_tick_params(color=MUTED, labelcolor=MUTED)
    fig.subplots_adjust(left=0.14, right=0.94, bottom=0.24, top=0.84)
    head(fig, "02 — Correlation matrix", "Pearson r across numeric features · strong pairs drive the insight panel")
    save(fig, "02-correlation.png")
    return corr


def signal(df, corr):
    target = "price"
    ranked = corr[target].drop(target).abs().sort_values()
    fig, ax = plt.subplots(figsize=(10, 5.4))
    colors = [ACCENT if v >= 0.4 else ACCENT2 for v in ranked.values]
    bars = ax.barh(ranked.index, ranked.values, color=colors, alpha=0.88, height=0.62)
    for bar, value in zip(bars, ranked.values):
        ax.text(value + 0.012, bar.get_y() + bar.get_height() / 2, f"{value:.2f}",
                va="center", color="#dbe6f7", fontsize=10)
    ax.set_xlim(0, max(ranked.values) * 1.16)
    ax.set_xlabel("|Pearson r| with price")
    ax.grid(axis="x")
    ax.set_axisbelow(True)
    fig.subplots_adjust(left=0.15, right=0.95, bottom=0.16, top=0.83)
    head(fig, "03 — Signal scan", "Which numeric features move with price · ranking used for smart insights")
    save(fig, "03-signal-scan.png")


def categories(df):
    columns = ["furnishingstatus", "mainroad", "airconditioning", "guestroom"]
    fig, axes = plt.subplots(2, 2, figsize=(12, 6.6))
    colors = [ACCENT, ACCENT2, LIME, PINK]
    for ax, column, color in zip(axes.flat, columns, colors):
        counts = df[column].value_counts().sort_values()
        ax.barh(counts.index, counts.values, color=color, alpha=0.85, height=0.6)
        for value, label in zip(counts.values, counts.index):
            ax.text(value + 6, list(counts.index).index(label), f"{value} ({value / len(df):.0%})",
                    va="center", color=MUTED, fontsize=9)
        ax.set_xlim(0, counts.max() * 1.22)
        ax.set_title(column)
        ax.grid(axis="x")
        ax.set_axisbelow(True)
    fig.subplots_adjust(left=0.13, right=0.95, bottom=0.09, top=0.83, hspace=0.45, wspace=0.28)
    head(fig, "04 — Categorical composition", "Level frequencies for the boolean and nominal columns · 545 rows")
    save(fig, "04-categories.png")


def main():
    df = pd.read_csv(DATA)
    print(f"dataset: {df.shape[0]} rows x {df.shape[1]} columns")
    print(f"missing cells: {int(df.isna().sum().sum())}")
    print(f"duplicate rows: {int(df.duplicated().sum())}")

    corr = correlation(df)
    pairs = (
        corr.where(np.triu(np.ones(corr.shape), k=1).astype(bool))
        .stack()
        .sort_values(key=abs, ascending=False)
    )
    print("top pairs:", ", ".join(f"{a}<->{b}={v:.2f}" for (a, b), v in pairs.head(3).items()))

    distributions(df)
    signal(df, corr)
    categories(df)


if __name__ == "__main__":
    main()
