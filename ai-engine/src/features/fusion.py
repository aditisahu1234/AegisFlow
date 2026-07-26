import numpy as np
import pandas as pd


DEVIATION_COLUMN = "normality_deviation"
PERCENTILE_COLUMN = "normality_percentile"


def fuse_normality_features(
    X: pd.DataFrame,
    deviation_scores: np.ndarray,
    percentile_ranks: np.ndarray,
) -> pd.DataFrame:
    """
    Construct the Normality Fusion representation:

        x' = [x, s(x), r(x)]

    where:
        x    = original preprocessed features
        s(x) = Isolation Forest deviation score
        r(x) = ECDF percentile rank
    """

    if len(X) != len(deviation_scores):
        raise ValueError(
            "Feature matrix and deviation scores "
            "have different lengths."
        )

    if len(X) != len(percentile_ranks):
        raise ValueError(
            "Feature matrix and percentile ranks "
            "have different lengths."
        )

    deviation_scores = np.asarray(
        deviation_scores,
        dtype=np.float64,
    )

    percentile_ranks = np.asarray(
        percentile_ranks,
        dtype=np.float64,
    )

    if not np.all(np.isfinite(deviation_scores)):
        raise ValueError(
            "Deviation scores contain NaN or infinity."
        )

    if not np.all(np.isfinite(percentile_ranks)):
        raise ValueError(
            "Percentile ranks contain NaN or infinity."
        )

    if np.any(
        (percentile_ranks < 0)
        | (percentile_ranks > 1)
    ):
        raise ValueError(
            "Percentile ranks must lie in [0, 1]."
        )

    fused = X.copy()

    fused[DEVIATION_COLUMN] = deviation_scores
    fused[PERCENTILE_COLUMN] = percentile_ranks

    return fused