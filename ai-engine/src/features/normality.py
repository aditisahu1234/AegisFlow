import joblib
import numpy as np

from src.config import MODEL_ARTIFACTS_DIR


ECDF_REFERENCE_PATH = (
    MODEL_ARTIFACTS_DIR
    / "ecdf_reference.joblib"
)


def fit_ecdf_reference(
    deviation_scores: np.ndarray,
) -> np.ndarray:
    """
    Build the empirical reference distribution
    from TRAINING deviation scores only.
    """

    scores = np.asarray(
        deviation_scores,
        dtype=np.float64,
    )

    if scores.ndim != 1:
        raise ValueError(
            "Deviation scores must be one-dimensional."
        )

    if scores.size == 0:
        raise ValueError(
            "Cannot build ECDF from empty scores."
        )

    if not np.all(np.isfinite(scores)):
        raise ValueError(
            "Deviation scores contain NaN or infinity."
        )

    return np.sort(scores)


def transform_ecdf(
    deviation_scores: np.ndarray,
    reference_scores: np.ndarray,
) -> np.ndarray:
    """
    Convert deviation scores into empirical
    percentile ranks using the TRAINING reference.

    r(x) = (# reference scores <= s(x)) / N
    """

    scores = np.asarray(
        deviation_scores,
        dtype=np.float64,
    )

    reference = np.asarray(
        reference_scores,
        dtype=np.float64,
    )

    positions = np.searchsorted(
        reference,
        scores,
        side="right",
    )

    percentile_ranks = (
        positions / reference.size
    )

    return percentile_ranks


def save_ecdf_reference(
    reference_scores: np.ndarray,
) -> None:

    MODEL_ARTIFACTS_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    joblib.dump(
        reference_scores,
        ECDF_REFERENCE_PATH,
    )


def load_ecdf_reference() -> np.ndarray:

    return joblib.load(
        ECDF_REFERENCE_PATH
    )