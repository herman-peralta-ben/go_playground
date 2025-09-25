package com.herman.example.android_go_playground.composables.example

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.herman.example.android_go_playground.ui.theme.AndroidGoPlaygroundTheme

private const val EXAMPLE_START_DESTINATION = "menu"

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ExampleApp(
    title: String,
    examples: List<Example>,
) {
    val navController = rememberNavController()
    val currentTitle = remember { mutableStateOf(title) }
    val allScreens = flattenExamples("", examples)

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    val navBackStackEntry by navController.currentBackStackEntryAsState()
                    val currentRoute = navBackStackEntry?.destination?.route
                    val isRoot = currentRoute == EXAMPLE_START_DESTINATION

                    TopAppBar(
                        title = { Text(if (isRoot) title else currentTitle.value) }
                    )
                }
            )
        }
    ) { paddingValues ->
        NavHost(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues),
            navController = navController,
            startDestination = EXAMPLE_START_DESTINATION,
        ) {
            composable(EXAMPLE_START_DESTINATION) {
                ExampleMenu(currentTitle, examples, navController, "")
            }
            for ((route, screen) in allScreens) {
                composable(route) { screen() }
            }
        }
    }
}

@Composable
private fun ExampleMenu(
    title: MutableState<String>,
    examples: List<Example>,
    navController: NavController,
    parentPath: String
) {
    Column(modifier = Modifier.padding(16.dp)) {
        examples.forEach { example ->
            val route = (parentPath + "/" + example.label.lowercase()).replace(" ", "")
            Text(
                text = example.label,
                modifier = Modifier
                    .fillMaxWidth()
                    .clickable {
                        if (example.screen != null) {
                            title.value = example.label
                            navController.navigate(route)
                        }
                    }
                    .padding(vertical = 8.dp),
                style = MaterialTheme.typography.bodyLarge
            )

            if (example.subExamples.isNotEmpty()) {
                ExampleMenu(title, example.subExamples.values.toList(), navController, route)
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
private fun ExampleAppPreview() {
    AndroidGoPlaygroundTheme {
        ExampleApp(
            "Go Examples", listOf(
                Example(
                    label = "Examples",
                    subExamples = mapOf(
                        "examples_1" to Example(
                            label = "Examples1",
                            screen = { Text("Examples1") },
                        ),
                        "examples_2" to Example(
                            label = "Examples2",
                            subExamples = mapOf(
                                "examples_2/1" to Example(
                                    label = "Example2 / 1",
                                    screen = { Text("Example2 / 1") },
                                ),
                                "examples_2/2" to Example(
                                    label = "Example2 / 2",
                                    screen = { Text("Example2 / 2") },
                                ),
                            ),
                        ),
                        "examples_3" to Example(
                            label = "Examples3",
                            screen = { Text("Examples3") },
                        )
                    ),
                ),
            ),
        )
    }
}
