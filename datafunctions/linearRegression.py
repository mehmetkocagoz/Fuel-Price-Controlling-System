import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn import metrics
import numpy as np

dataFrame = pd.read_csv("price_data.csv")
print(dataFrame.corr(method="pearson",numeric_only=True))
x = dataFrame[['brentoilprice','exchange_column']]
y = dataFrame['fuelprice']
X_train, X_test, y_train, y_test = train_test_split(x, y, test_size=0.2, random_state=101)
lm = LinearRegression()
lm.fit(X_train,y_train)
coeff_df = pd.DataFrame(lm.coef_,x.columns,columns=['Coefficient'])
print(coeff_df)
predictions = lm.predict(X_test)
plt.scatter(y_test,predictions)
plt.show()
MAE = metrics.mean_absolute_error(y_test,predictions)
MSE = metrics.mean_squared_error(y_test,predictions)
RMSE = np.sqrt(metrics.mean_squared_error(y_test,predictions)) 
print('MAE:',MAE)
print('MSE:',MSE)
print('RMSE:',RMSE)