package com.herman.example.android_go_playground

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import com.herman.example.android_go_playground.composables.example.Example
import com.herman.example.android_go_playground.composables.example.ExampleApp
import com.herman.example.android_go_playground.ui.theme.AndroidGoPlaygroundTheme

private val examples = listOf(
    Example(
        label = "Examples",
        subExamples = mapOf(
            "examples_1" to Example(
                label = "Examples1",
                screen = { ExampleScreen("Examples1") },
            ),
            "examples_2" to Example(
                label = "Examples2",
                subExamples = mapOf(
                    "examples_2/1" to Example(
                        label = "Example2 / 1",
                        screen = { ExampleScreen("Example2 / 1") },
                    ),
                    "examples_2/2" to Example(
                        label = "Example2 / 2",
                        screen = { ExampleScreen("Example2 / 2") },
                    ),
                ),
            ),
            "examples_3" to Example(
                label = "Examples3",
                subExamples = mapOf(
                    "examples_3/1" to Example(
                        label = "Example3 / 1",
                        screen = { ExampleScreen("Example3 / 1") },
                    ),
                ),
            ),
        ),
    ),
)

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            AndroidGoPlaygroundTheme {
                ExampleApp("Go Examples", examples)
            }
        }
    }
}

@Composable
private fun ExampleScreen(name: String) {
    Column(verticalArrangement = Arrangement.Center) {
        Text(name)
    }
}

@Preview(showBackground = true)
@Composable
private fun ExampleScreenPreview() {
    AndroidGoPlaygroundTheme {
        ExampleScreen("Android")
    }
}
