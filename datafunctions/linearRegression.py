import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn import metrics
import numpy as np

dataFrame = pd.read_csv("data.csv")
print(dataFrame.corr(method="pearson",numeric_only=True))
x = dataFrame[['BrentOilPrice','USD/TRY']]
y = dataFrame['FuelPrice']
X_train, X_test, y_train, y_test = train_test_split(x, y, test_size=0.2, random_state=101)
lm = LinearRegression()
lm.fit(X_train,y_train)
coeff_df = pd.DataFrame(lm.coef_,x.columns,columns=['Coefficient'])
print(coeff_df)
predictions = lm.predict(X_test)
#plt.scatter(y_test,predictions)
#plt.show()
MAE = metrics.mean_absolute_error(y_test,predictions)
MSE = metrics.mean_squared_error(y_test,predictions)
RMSE = np.sqrt(metrics.mean_squared_error(y_test,predictions)) 
print('MAE:',MAE)
print('MSE:',MSE)
print('RMSE:',RMSE)


dataFrame2 = pd.read_csv("cleanData.csv")
print(dataFrame.corr(method="pearson",numeric_only=True))
x2 = dataFrame2[['BrentOilPrice','USD/TRY']]
y2 = dataFrame2['FuelPrice']
X2_train, X2_test, y2_train, y2_test = train_test_split(x2, y2, test_size=0.2, random_state=101)
lm2 = LinearRegression()
lm2.fit(X2_train,y2_train)
coeff_df = pd.DataFrame(lm2.coef_,x2.columns,columns=['Coefficient'])
print(coeff_df)
predictions2 = lm2.predict(X2_test)
#plt.scatter(y2_test,predictions)
#plt.show()
MAE2 = metrics.mean_absolute_error(y2_test,predictions2)
MSE2 = metrics.mean_squared_error(y2_test,predictions2)
RMSE2 = np.sqrt(metrics.mean_squared_error(y2_test,predictions2)) 
print('MAE:',MAE2)
print('MSE:',MSE2)
print('RMSE:',RMSE2)


'''
By doing same operations on clean data, we can see that the results are better.
MAE, MSE and RMSE values are lower than the previous ones.
It means that the model is more accurate.
'''
