import androidx.activity.ComponentActivity
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.platform.LocalContext
import com.stripe.android.googlepaylauncher.GooglePayLauncher

@Composable
fun GooglePayButton(
    amount: Int,
    currencyCode: String,
    onResult: (Result<String>) -> Unit
) {
    val context = LocalContext.current
    val activity = context as? ComponentActivity
        ?: throw IllegalStateException("GooglePayButton debe usarse dentro de una Activity")

    var launcher: GooglePayLauncher? by remember { mutableStateOf(null) }

    LaunchedEffect(Unit) {
        launcher = GooglePayLauncher(
            activity = activity,
            config = GooglePayLauncher.Config(
                environment = GooglePayLauncher.Config.Environment.Test,
                merchantCountryCode = "MX",
                merchantName = "Mi App"
            ),
            readyCallback = { isReady ->
                if (isReady) {
                    launcher?.present(
                        GooglePayLauncher.PaymentDataArgs(
                            currencyCode = currencyCode,
                            amount = amount
                        )
                    )
                } else {
                    onResult(Result.failure(Exception("Google Pay no está disponible")))
                }
            },
            resultCallback = { result ->
                when (result) {
                    is GooglePayLauncher.Result.Completed -> {
                        val paymentMethodId = result.paymentMethod?.id
                        if (paymentMethodId != null) {
                            onResult(Result.success(paymentMethodId))
                        } else {
                            onResult(Result.failure(Exception("No se recibió paymentMethod")))
                        }
                    }
                    is GooglePayLauncher.Result.Canceled -> {
                        onResult(Result.failure(Exception("Cancelado por el usuario")))
                    }
                    is GooglePayLauncher.Result.Failed -> {
                        onResult(Result.failure(result.error))
                    }
                }
            }
        )
    }

    Button(onClick = {
        launcher?.present(
            GooglePayLauncher.PaymentDataArgs(
                currencyCode = currencyCode,
                amount = amount
            )
        )
    }) {
        Text("Pagar con Google Pay")
    }
}
